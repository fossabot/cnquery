package vsphere

import (
	"errors"

	"github.com/rs/zerolog/log"
	"go.mondoo.io/mondoo/motor/asset"
	"go.mondoo.io/mondoo/motor/discovery/credentials"
	"go.mondoo.io/mondoo/motor/providers"
	"go.mondoo.io/mondoo/motor/providers/resolver"
	"go.mondoo.io/mondoo/motor/providers/vmwareguestapi"
	"go.mondoo.io/mondoo/motor/vault"
)

type VMGuestResolver struct{}

func (k *VMGuestResolver) Name() string {
	return "VmWare vSphere VM Guest Resolver"
}

func (r *VMGuestResolver) AvailableDiscoveryTargets() []string {
	return []string{}
}

func (k *VMGuestResolver) Resolve(root *asset.Asset, tc *providers.TransportConfig, cfn credentials.CredentialFn, sfn credentials.QuerySecretFn, userIdDetectors ...providers.PlatformIdDetector) ([]*asset.Asset, error) {
	resolved := []*asset.Asset{}

	// we leverage the vpshere transport to establish a connection
	m, err := resolver.NewMotorConnection(tc, cfn)
	if err != nil {
		return nil, err
	}
	defer m.Close()

	trans, ok := m.Transport.(*vmwareguestapi.Transport)
	if !ok {
		return nil, errors.New("could not initialize vsphere guest transport")
	}

	client := trans.Client()
	discoveryClient := New(client)

	// resolve vms
	vms, err := discoveryClient.ListVirtualMachines(tc)
	if err != nil {
		return nil, err
	}

	// add transport config for each vm
	for i := range vms {
		vm := vms[i]
		resolved = append(resolved, vm)
	}

	// filter the vms by inventoryPath
	inventoryPaths := []string{}
	inventoryPathFilter, ok := tc.Options["inventoryPath"]
	if ok {
		inventoryPaths = []string{inventoryPathFilter}
	}

	resolved = filter(resolved, func(a *asset.Asset) bool {
		inventoryPathLabel := a.Labels["vsphere.vmware.com/inventory-path"]

		return contains(inventoryPaths, inventoryPathLabel)
	})

	if len(resolved) == 1 {
		a := resolved[0]
		a.Connections = []*providers.TransportConfig{tc}

		// find the secret reference for the asset
		EnrichVsphereToolsConnWithSecrets(a, cfn, sfn)

		return []*asset.Asset{a}, nil
	} else {
		return nil, errors.New("could not resolve vSphere vm")
	}
}

func EnrichVsphereToolsConnWithSecrets(a *asset.Asset, cfn credentials.CredentialFn, sfn credentials.QuerySecretFn) {
	// search secret for vm
	// NOTE: we do not use `common.EnrichAssetWithSecrets(a, sfn)` here since vmware requires two secrets at the same time
	for j := range a.Connections {
		conn := a.Connections[j]

		// special handling for vsphere vm config
		if conn.Backend == providers.TransportBackend_CONNECTION_VSPHERE_VM {
			var creds *vault.Credential

			secretRefCred, err := sfn(a)
			if err == nil && secretRefCred != nil {
				creds, err = cfn(secretRefCred)
			}

			if err == nil && creds != nil {
				if conn.Options == nil {
					conn.Options = map[string]string{}
				}
				conn.Options["guestUser"] = creds.User
				conn.Options["guestPassword"] = string(creds.Secret)
			}
		} else {
			log.Warn().Str("name", a.Name).Msg("could not determine credentials for asset")
		}
	}
}
