// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package provider

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/providers-sdk/v1/inventory"
	"go.mondoo.com/cnquery/providers/os/connection/shared"
	"go.mondoo.com/cnquery/providers/os/detector"
	"go.mondoo.com/cnquery/providers/os/id/awsec2"
	"go.mondoo.com/cnquery/providers/os/id/awsecs"
	"go.mondoo.com/cnquery/providers/os/id/clouddetect"
	"go.mondoo.com/cnquery/providers/os/id/hostname"
	"go.mondoo.com/cnquery/providers/os/id/ids"
	"go.mondoo.com/cnquery/providers/os/id/machineid"
	"go.mondoo.com/cnquery/providers/os/id/sshhostkey"
)

type PlatformFingerprint struct {
	PlatformIDs   []string
	Name          string
	Runtime       string
	Kind          string
	RelatedAssets []PlatformFingerprint
}

type PlatformInfo struct {
	IDs                []string
	Name               string
	RelatedPlatformIDs []string
}

func IdentifyPlatform(conn shared.Connection, p *inventory.Platform, idDetectors []string) (*PlatformFingerprint, error) {
	var ok bool
	if p == nil {
		p, ok = detector.DetectOS(conn)
		if !ok {
			return nil, errors.New("cannot detect os")
		}
	}

	var fingerprint PlatformFingerprint
	var ids []string
	var relatedIds []string

	for i := range idDetectors {
		idDetector := idDetectors[i]
		platformInfo, err := GatherPlatformInfo(conn, p, idDetector)
		if err != nil {
			// we only err if we found zero platform ids, if we try multiple, a fail of an individual one is okay
			log.Debug().Err(err).Str("detector", string(idDetector)).Msg("could not determine platform info")
			continue
		}
		if len(platformInfo.IDs) > 0 {
			ids = append(ids, platformInfo.IDs...)
		}
		if len(platformInfo.RelatedPlatformIDs) > 0 {
			relatedIds = append(relatedIds, platformInfo.RelatedPlatformIDs...)
		}

		if len(platformInfo.Name) > 0 {
			fingerprint.Name = platformInfo.Name
		} else {
			// check if we get a name for the asset, eg. aws instance id
			for _, id := range platformInfo.IDs {
				name := GatherNameForPlatformId(id)
				if name != "" {
					fingerprint.Name = name
				}
			}
		}
		// check whether we can extract runtime and kind information
		for _, id := range platformInfo.IDs {
			runtime, kind := ExtractPlatformAndKindFromPlatformId(id)
			if runtime != "" {
				fingerprint.Runtime = runtime
				fingerprint.Kind = kind
			}
		}
	}

	// if we found zero platform ids something went wrong
	if len(ids) == 0 {
		return nil, errors.New("could not determine a platform identifier")
	}

	fingerprint.PlatformIDs = ids
	for _, v := range relatedIds {
		fingerprint.RelatedAssets = append(fingerprint.RelatedAssets, PlatformFingerprint{
			PlatformIDs: []string{v},
			Name:        GatherNameForPlatformId(v),
		})
	}

	log.Debug().Interface("id-detector", idDetectors).Strs("platform-ids", ids).Msg("detected platform ids")
	return &fingerprint, nil
}

func GatherNameForPlatformId(id string) string {
	if awsec2.IsValidMondooInstanceId(id) {
		structId, _ := awsec2.ParseMondooInstanceID(id)
		return structId.Id
	} else if accountID, err := awsec2.ParseMondooAccountID(id); err == nil {
		return fmt.Sprintf("AWS Account %s", accountID)
	}
	return ""
}

func ExtractPlatformAndKindFromPlatformId(id string) (string, string) {
	if awsec2.ParseEc2PlatformID(id) != nil {
		return "aws-ec2", "virtual-machine"
	} else if awsec2.IsValidMondooAccountId(id) {
		return "aws", "api"
	} else if awsecs.IsValidMondooECSContainerId(id) {
		return "aws-ecs", "container"
	}
	return "", ""
}

func GatherPlatformInfo(conn shared.Connection, pf *inventory.Platform, idDetector string) (*PlatformInfo, error) {
	var identifier string
	switch {
	case idDetector == ids.IdDetector_Hostname:
		// NOTE: we need to be careful with hostname's since they are not required to be unique
		hostname, ok := hostname.Hostname(conn, pf)
		if ok && len(hostname) > 0 {
			identifier = "//platformid.api.mondoo.app/hostname/" + hostname
			return &PlatformInfo{
				IDs:                []string{identifier},
				Name:               hostname,
				RelatedPlatformIDs: []string{},
			}, nil
		}
		return &PlatformInfo{}, nil
	case idDetector == ids.IdDetector_MachineID:
		guid, hostErr := machineid.MachineId(conn, pf)
		if hostErr == nil && len(guid) > 0 {
			identifier = "//platformid.api.mondoo.app/machineid/" + guid
			return &PlatformInfo{
				IDs:                []string{identifier},
				Name:               "",
				RelatedPlatformIDs: []string{},
			}, hostErr
		}
		return &PlatformInfo{}, nil
	case idDetector == ids.IdDetector_AwsEcs:
		metadata, err := awsecs.Resolve(conn, pf)
		if err != nil {
			return nil, err
		}
		ident, err := metadata.Identify()
		if err != nil {
			return nil, err
		}
		if len(ident.PlatformIds) != 0 {
			return &PlatformInfo{
				IDs:                ident.PlatformIds,
				Name:               ident.Name,
				RelatedPlatformIDs: []string{ident.AccountPlatformID},
			}, nil
		}
		return &PlatformInfo{}, nil
	case idDetector == ids.IdDetector_CloudDetect:
		identifier, name, relatedIdentifiers := clouddetect.Detect(conn, pf)
		if identifier != "" {
			return &PlatformInfo{
				IDs:                []string{identifier},
				Name:               name,
				RelatedPlatformIDs: relatedIdentifiers,
			}, nil
		}
		return &PlatformInfo{}, nil
	case idDetector == ids.IdDetector_SshHostkey:
		identifier, err := sshhostkey.Detect(conn, pf)
		if err != nil {
			return nil, err
		}
		return &PlatformInfo{
			IDs:                identifier,
			Name:               "",
			RelatedPlatformIDs: []string{},
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("the provided id-detector is not supported: %s", idDetector))
	}
}
