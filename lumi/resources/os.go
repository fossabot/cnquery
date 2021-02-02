package resources

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go.mondoo.io/mondoo/lumi/resources/reboot"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mondoo.io/mondoo/llx"
	"go.mondoo.io/mondoo/lumi/resources/packages"
	"go.mondoo.io/mondoo/lumi/resources/platformid"
	"go.mondoo.io/mondoo/lumi/resources/uptime"
	"go.mondoo.io/mondoo/lumi/resources/windows"
	"go.mondoo.io/mondoo/motor/motorid/hostname"
)

func (p *lumiOs) id() (string, error) {
	return "os", nil
}

func (p *lumiOs) GetRebootpending() (interface{}, error) {
	// try to collect if a reboot is required, fails for static images
	rb, err := reboot.New(p.Runtime.Motor)
	if err != nil {
		return nil, err
	}
	return rb.RebootPending()
}

func (p *lumiOs) getUnixEnv() (map[string]interface{}, error) {
	rawCmd, err := p.Runtime.CreateResource("command", "command", "env")
	if err != nil {
		return nil, err
	}
	cmd := rawCmd.(Command)

	out, err := cmd.Stdout()
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{}
	lines := strings.Split(out, "\n")
	for i := range lines {
		parts := strings.SplitN(lines[i], "=", 2)
		if len(parts) != 2 {
			continue
		}
		res[parts[0]] = parts[1]
	}

	return res, nil
}

func (p *lumiOs) getWindowsEnv() (map[string]interface{}, error) {
	rawCmd, err := p.Runtime.CreateResource("powershell",
		"script", "Get-ChildItem Env:* | ConvertTo-Json",
	)
	if err != nil {
		return nil, err
	}
	cmd := rawCmd.(Powershell)

	out, err := cmd.Stdout()
	if err != nil {
		return nil, err
	}

	return windows.ParseEnv(strings.NewReader(out))
}

func (p *lumiOs) GetEnv() (map[string]interface{}, error) {
	pf, err := p.Runtime.Motor.Platform()
	if err != nil {
		return nil, err
	}

	if pf.IsFamily("windows") {
		return p.getWindowsEnv()
	}

	return p.getUnixEnv()
}

func (p *lumiOs) GetPath() ([]interface{}, error) {
	env, err := p.Env()
	if err != nil {
		return nil, err
	}

	rawPath, ok := env["PATH"]
	if !ok {
		return []interface{}{}, nil
	}

	path := rawPath.(string)
	parts := strings.Split(path, ":")
	res := make([]interface{}, len(parts))
	for i := range parts {
		res[i] = parts[i]
	}

	return res, nil
}

// returns uptime in nanoseconds
func (p *lumiOs) GetUptime() (*time.Time, error) {

	uptime, err := uptime.New(p.Runtime.Motor)
	if err != nil {
		return LumiTime(llx.DurationToTime(0)), err
	}

	t, err := uptime.Duration()
	if err != nil {
		return LumiTime(llx.DurationToTime(0)), err
	}

	// we get nano seconds but duration to time only takes seconds
	bootTime := time.Now().Add(-t)
	up := time.Now().Unix() - bootTime.Unix()
	return LumiTime(llx.DurationToTime(up)), nil
}

// func (p *lumiOs) GetRebootpending() ([]interface{}, error) {
// 	return nil, errors.New("not implemented")
// }

func (p *lumiOsUpdate) id() (string, error) {
	name, _ := p.Name()
	return name, nil
}

func (p *lumiOs) GetUpdates() ([]interface{}, error) {
	// find suitable system updates
	um, err := packages.ResolveSystemUpdateManager(p.Runtime.Motor)
	if um == nil || err != nil {
		return nil, fmt.Errorf("could not detect suiteable update manager for platform")
	}

	// retrieve all system updates
	updates, err := um.List()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve updates list for platform")
	}

	// create lumi update resources for each update
	osupdates := make([]interface{}, len(updates))
	log.Debug().Int("updates", len(updates)).Msg("lumi[updates]> found system updates")
	for i, update := range updates {

		lumiOsUpdate, err := p.Runtime.CreateResource("os.update",
			"name", update.Name,
			"severity", update.Severity,
			"category", update.Category,
			"restart", update.Restart,
			"format", um.Format(),
		)
		if err != nil {
			return nil, err
		}

		osupdates[i] = lumiOsUpdate.(OsUpdate)
	}

	// return the packages as new entries
	return osupdates, nil
}

func (s *lumiOs) GetHostname() (string, error) {
	platform, err := s.Runtime.Motor.Platform()
	if err != nil {
		return "", errors.New("cannot determine platform uuid")
	}

	return hostname.Hostname(s.Runtime.Motor.Transport, platform)
}

// returns the OS native machine UUID/GUID
func (s *lumiOs) GetMachineid() (string, error) {
	platform, err := s.Runtime.Motor.Platform()
	if err != nil {
		return "", errors.New("cannot determine platform uuid")
	}

	var uuidProvider platformid.UniquePlatformIDProvider
	for i := range platform.Family {
		if platform.Family[i] == "linux" {
			uuidProvider = &platformid.LinuxIdProvider{Motor: s.Runtime.Motor}
		}
	}

	if uuidProvider == nil && platform.Name == "macos" {
		uuidProvider = &platformid.MacOSIdProvider{Motor: s.Runtime.Motor}
	}

	if uuidProvider == nil {
		return "", errors.New("cannot determine platform uuid for " + platform.Name)
	}

	id, err := uuidProvider.ID()
	if err != nil {
		return "", errors.New("cannot determine platform uuid on known system " + platform.Name)
	}

	// TODO: we may want to inject that during compile time
	return HashedMachineID("3zXPqBRdu2zyspzplk7gxi1LEveYBrY0hdgCYv4M", id), nil
}

// We use a mechanism established by https://github.com/denisbrodbeck/machineid to
// derive the platform id in a reliable manner but we are not exposing the machine secret
func HashedMachineID(secret, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(secret))
	return hex.EncodeToString(mac.Sum(nil))
}
