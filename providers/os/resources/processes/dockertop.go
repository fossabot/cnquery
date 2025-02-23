// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package processes

import (
	"context"
	"fmt"
	"strconv"

	"go.mondoo.com/cnquery/providers/os/connection"
	"go.mondoo.com/cnquery/providers/os/connection/shared"
)

type DockerTopManager struct {
	conn shared.Connection
}

func (lpm *DockerTopManager) Name() string {
	return "Docker Top Process Manager"
}

// List lists the processes running in a Docker container. Note that currently this function returns child
// processes as well.
func (lpm *DockerTopManager) List() ([]*OSProcess, error) {
	dockerConn, ok := lpm.conn.(*connection.DockerContainerConnection)
	if !ok {
		return nil, fmt.Errorf("wrong transport type")
	}

	ctx := context.Background()
	client := dockerConn.Client

	// The Docker API uses ps underneath so we can provide any ps arguments we want here.
	resp, err := client.ContainerTop(ctx, dockerConn.ContainerId(), []string{"-o", "pid,user,comm,s,command"})
	if err != nil {
		return nil, err
	}

	// The docker API returns a list of strings for each process with the following format:
	// [0]: PID
	// [1]: USER
	// [2]: executable
	// [3]: state
	// [4]: command
	var procs []*OSProcess
	for _, p := range resp.Processes {
		pid, err := strconv.Atoi(p[0])
		if err != nil {
			continue
		}
		procs = append(procs, &OSProcess{
			Pid:          int64(pid), // This will be the PID inside the container
			Executable:   p[2],
			Command:      p[4],
			State:        p[3],
			SocketInodes: nil,
		})
	}

	return procs, nil
}

// check that the pid directory exists
func (lpm *DockerTopManager) Exists(pid int64) (bool, error) {
	procs, err := lpm.List()
	if err != nil {
		return false, err
	}

	for _, p := range procs {
		if p.Pid == pid {
			return true, nil
		}
	}

	return false, nil
}

func (lpm *DockerTopManager) Process(pid int64) (*OSProcess, error) {
	procs, err := lpm.List()
	if err != nil {
		return nil, err
	}

	for _, p := range procs {
		if p.Pid == pid {
			return p, nil
		}
	}
	return nil, fmt.Errorf("process with PID %d does not exist", pid)
}
