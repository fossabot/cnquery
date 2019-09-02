package events

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"go.mondoo.io/mondoo/motor/motoros/types"
)

// the job state
type JobState int32

const (
	// pending is the default state
	Job_PENDING    JobState = 0
	Job_RUNNING    JobState = 1
	Job_TERMINATED JobState = 2
)

type Job struct {
	ID string

	Runnable func(types.Transport) (types.Observable, error)
	Callback []func(types.Observable)

	State        JobState
	ScheduledFor time.Time
	Interval     time.Duration
	// -1 means infinity
	Repeat int32

	Metrics struct {
		RunAt     time.Time
		Duration  time.Duration
		Count     int32
		Errors    int32
		Successes int32
	}
}

func (j *Job) sanitize() error {
	// ensure we have an id
	if len(j.ID) == 0 {
		j.ID = uuid.Must(uuid.NewV4()).String()
	}

	// verify that the interval is set for the job, otherwise overwrite with the default
	if j.Interval == 0 {
		j.Interval = time.Duration(60 * time.Second)
	}

	// verify that we have the required things for a schedule
	if j.ScheduledFor.Before(time.Now().Add(time.Duration(-10 * time.Second))) {
		return errors.New("schedule for the past")
	}

	if j.Runnable == nil {
		return errors.New("no runnable defined")
	}

	if len(j.Callback) == 0 {
		return errors.New("no callback defined")
	}

	return nil
}

func (j *Job) SetInfinity() {
	j.Repeat = -1
}

func (j *Job) isPending() bool {
	return j.State == Job_PENDING
}

func NewJobManager(transport types.Transport) *JobManager {
	jm := &JobManager{transport: transport}

	// stores all callbacks for by subscriber
	jm.jobs = make(map[string]*Job)
	jm.jobMutex = &sync.Mutex{}

	jm.Serve()
	return jm
}

type JobManagerMetrics struct {
	Jobs int
}

type JobManager struct {
	transport  types.Transport
	quit       chan bool
	jobs       map[string]*Job
	jobMutex   *sync.Mutex
	jobMetrics JobManagerMetrics
}

// Schedule stores the job in the run list and sanitize the job before execution
func (jm *JobManager) Schedule(job *Job) (string, error) {
	// ensure all defaults are set
	err := job.sanitize()
	if err != nil {
		return "", err
	}

	log.Debug().Str("jobid", job.ID).Msg("motor.job> schedule new job")

	// store job, with a mutex
	jm.jobMutex.Lock()
	jm.jobs[job.ID] = job
	jm.jobMutex.Unlock()

	// return job id
	return job.ID, nil
}

func (jm *JobManager) GetJob(jobid string) (*Job, error) {
	jm.jobMutex.Lock()
	job, ok := jm.jobs[jobid]
	jm.jobMutex.Unlock()
	if !ok {
		return nil, errors.New("job " + jobid + " does not exist")
	}

	return job, nil
}

func (jm *JobManager) Delete(jobid string) error {
	log.Debug().Str("jobid", jobid).Msg("motor.job> delete job")
	jm.jobMutex.Lock()
	delete(jm.jobs, jobid)
	jm.jobMutex.Unlock()
	return nil
}

func (jm *JobManager) Metrics() *JobManagerMetrics {
	jm.jobMetrics.Jobs = len(jm.jobs)
	return &jm.jobMetrics
}

// Serve creates a goroutine and runs jobs in the background
func (jm *JobManager) Serve() {
	// create a new channel and starte a new go routine
	jm.quit = make(chan bool)
	go func() {
		for {
			select {
			case <-jm.quit:
				return
			default:
				// fetch job
				job, err := jm.nextJob()

				if err == nil {
					// run job
					jm.Run(job)

					// if repeat is 0 and it is not the last iteration of a reoccuring task,
					// we need to remove the job
					if job.Repeat == 0 && job.State == Job_TERMINATED {
						jm.Delete(job.ID)
					}
				}

				// TODO: wake up, when new jobs come in
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (jm *JobManager) Run(job *Job) error {
	log.Debug().Str("jobid", job.ID).Msg("motor.job> run job")
	job.Metrics.RunAt = time.Now()

	// execute job
	observable, err := job.Runnable(jm.transport)

	// update metrics
	job.Metrics.Count = job.Metrics.Count + 1
	if err != nil {
		job.Metrics.Errors = job.Metrics.Errors + 1
	} else {
		job.Metrics.Successes = job.Metrics.Successes + 1
	}

	// determine the next run or delete the job
	if job.Repeat != 0 {
		job.ScheduledFor = time.Now().Add(job.Interval)
		log.Debug().Str("jobid", job.ID).Time("time", job.ScheduledFor).Msg("motor.job> scheduled job for the future")
		job.State = Job_PENDING
	} else {
		log.Debug().Str("jobid", job.ID).Msg("motor.job> last run for this job, yeah")
		job.State = Job_TERMINATED
	}

	// if we have a positive repeat, we need to decrement
	if job.Repeat > 0 {
		job.Repeat = job.Repeat - 1
	}

	// calc duration
	job.Metrics.Duration = time.Now().Sub(job.Metrics.RunAt)
	log.Debug().Str("jobid", job.ID).Msg("motor.job> completed")

	// send observable to all subscribers
	// since this call is synchronous in the same go routine, we need to do this as the last step, to ensure
	// all job planning is completed before a potential canceling comes in.
	log.Debug().Str("jobid", job.ID).Msg("motor.job> call subscriber")
	for _, subscriber := range job.Callback {
		subscriber(observable)
	}

	return nil
}

// nextJob looks for the oldest job and does that one first
func (jm *JobManager) nextJob() (*Job, error) {
	// use lock to prevent concurrent access on that list
	var oldestJob *Job
	oldest := time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)

	// iterate over list of jobs of pending jobs and find the oldest one
	jm.jobMutex.Lock()
	now := time.Now()
	for _, job := range jm.jobs {
		if job.State == Job_PENDING && oldest.After(job.ScheduledFor) && job.ScheduledFor.Before(now) {
			oldest = job.ScheduledFor
			oldestJob = job
		}
	}

	// set the job to running to ensure other parallel go routines do not fetch
	// the same job
	if oldestJob != nil {
		oldestJob.State = Job_RUNNING
	}

	jm.jobMutex.Unlock()

	if oldestJob == nil {
		return nil, errors.New("no job available")
	}

	// extrats the next run from the nextruns
	return oldestJob, nil
}

// TeadDown deletes all
func (jm *JobManager) TearDown() {
	// ensures the go routines are canceled
	jm.quit <- true
}
