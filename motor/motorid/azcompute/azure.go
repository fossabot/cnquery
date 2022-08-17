package azcompute

import (
	"fmt"

	"go.mondoo.io/mondoo/motor/providers/os"

	"github.com/pkg/errors"
	"go.mondoo.io/mondoo/motor/platform"
)

type instanceMetadata struct {
	Compute struct {
		ResourceID     string `json:"resourceID"`
		SubscriptionID string `json:"subscriptionId"`
		Tags           string `json:"tags"`
	} `json:"compute"`
}

type InstanceIdentifier interface {
	InstanceID() (string, error)
}

func Resolve(provider os.OperatingSystemProvider, pf *platform.Platform) (InstanceIdentifier, error) {
	if pf.IsFamily(platform.FAMILY_UNIX) || pf.IsFamily(platform.FAMILY_WINDOWS) {
		return NewCommandInstanceMetadata(provider, pf), nil
	}
	return nil, errors.New(fmt.Sprintf("azure compute id detector is not supported for your asset: %s %s", pf.Name, pf.Version))
}
