package platformid

import (
	"io/ioutil"
	"strings"

	"go.mondoo.io/mondoo/motor"
)

// LinuxIdProvider read the following files to extract the machine id
// "/var/lib/dbus/machine-id" and "/etc/machine-id"
// TODO: this approach is only reliable for systemd managed machines
type LinuxIdProvider struct {
	Motor *motor.Motor
}

func (p *LinuxIdProvider) Name() string {
	return "Linux Machine ID"
}

func (p *LinuxIdProvider) ID() (string, error) {
	content, err := p.retrieveFile("/var/lib/dbus/machine-id")
	if err != nil {
		content, err = p.retrieveFile("/etc/machine-id")
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(strings.ToLower(string(content))), nil
}

func (p *LinuxIdProvider) retrieveFile(path string) ([]byte, error) {
	f, err := p.Motor.Transport.File(path)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return content, nil
}
