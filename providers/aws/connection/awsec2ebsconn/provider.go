// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsec2ebsconn

import (
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"go.mondoo.com/cnquery/providers-sdk/v1/inventory"
	awsec2ebstypes "go.mondoo.com/cnquery/providers/aws/connection/awsec2ebsconn/types"
	"go.mondoo.com/cnquery/providers/os/connection"
	"go.mondoo.com/cnquery/providers/os/connection/shared"
	"go.mondoo.com/cnquery/providers/os/connection/snapshot"
	"go.mondoo.com/cnquery/providers/os/detector"
)

const (
	EBSConnectionType shared.ConnectionType = "ebs"
)

type AwsEbsConnection struct {
	id                  uint32
	asset               *inventory.Asset
	FsProvider          *connection.FileSystemConnection
	scannerRegionEc2svc *ec2.Client
	targetRegionEc2svc  *ec2.Client
	config              aws.Config
	opts                map[string]string
	scannerInstance     *awsec2ebstypes.InstanceId // the instance the transport is running on
	scanVolumeInfo      *awsec2ebstypes.VolumeInfo // the info of the volume we attached to the instance
	target              awsec2ebstypes.TargetInfo  // info about the target
	targetType          string                     // the type of object we're targeting (instance, volume, snapshot)
	volumeMounter       *snapshot.VolumeMounter
}

// New creates a new aws-ec2-ebs provider
// It expects to be running on an ec2 instance with ssm iam role and
// permissions for copy snapshot, create snapshot, create volume, attach volume, detach volume
func NewAwsEbsConnection(id uint32, conf *inventory.Config, asset *inventory.Asset) (*AwsEbsConnection, error) {
	log.Debug().Msg("new aws ebs connection")
	// TODO: validate the expected permissions here
	// TODO: allow custom aws config
	// 1. validate; load
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "could not load aws configuration")
	}
	i, err := RawInstanceInfo(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not load instance info: aws-ec2-ebs provider only valid on ec2 instances")
	}

	// ec2 client for the scanner region
	cfg.Region = i.Region
	scannerSvc := ec2.NewFromConfig(cfg)

	// ec2 client for the target region
	cfgCopy := cfg.Copy()
	cfgCopy.Region = conf.Options["region"]
	targetSvc := ec2.NewFromConfig(cfgCopy)

	shell := []string{"sh", "-c"}
	volumeMounter := snapshot.NewVolumeMounter(shell)

	// 2. create provider instance
	c := &AwsEbsConnection{
		config: cfg,
		opts:   conf.Options,
		target: awsec2ebstypes.TargetInfo{
			PlatformId: conf.PlatformId,
			Region:     conf.Options["region"],
			Id:         conf.Options["id"],
		},
		targetType: conf.Options["type"],
		scannerInstance: &awsec2ebstypes.InstanceId{
			Id:      i.InstanceID,
			Region:  i.Region,
			Account: i.AccountID,
			Zone:    i.AvailabilityZone,
		},
		targetRegionEc2svc:  targetSvc,
		scannerRegionEc2svc: scannerSvc,
		volumeMounter:       volumeMounter,
		asset:               asset,
	}
	log.Debug().Interface("info", c.target).Str("type", c.targetType).Msg("target")

	ctx := context.Background()

	// 3. validate
	instanceinfo, volumeid, snapshotid, err := c.Validate(ctx)
	if err != nil {
		return c, errors.Wrap(err, "unable to validate")
	}

	// 4. setup the volume for scanning
	// check if we got the no setup override option. this implies the target volume is already attached to the instance
	// this is used in cases where we need to test a snapshot created from a public marketplace image. the volume gets attached to a brand
	// new instance, and then that instance is started and we scan the attached fs
	if conf.Options[snapshot.NoSetup] == "true" {
		log.Info().Msg("skipping setup step")
	} else {
		var ok bool
		var err error
		switch c.targetType {
		case awsec2ebstypes.EBSTargetInstance:
			ok, err = c.SetupForTargetInstance(ctx, instanceinfo)
		case awsec2ebstypes.EBSTargetVolume:
			ok, err = c.SetupForTargetVolume(ctx, *volumeid)
		case awsec2ebstypes.EBSTargetSnapshot:
			ok, err = c.SetupForTargetSnapshot(ctx, *snapshotid)
		default:
			return c, errors.New("invalid target type")
		}
		if err != nil {
			log.Error().Err(err).Msg("unable to complete setup step")
			c.Close()
			return c, err
		}
		if !ok {
			return c, errors.New("something went wrong; unable to complete setup for ebs volume scan")
		}
	}

	// Mount Volume
	err = c.volumeMounter.Mount()
	if err != nil {
		log.Error().Err(err).Msg("unable to complete mount step")
		return c, err
	}
	// Create and initialize fs provider
	fsConn, err := connection.NewFileSystemConnection(id, &inventory.Config{
		Path:       c.volumeMounter.ScanDir,
		Backend:    "fs",
		PlatformId: conf.PlatformId,
		Options:    conf.Options,
	}, asset)
	if err != nil {
		return nil, err
	}
	c.FsProvider = fsConn
	var ok bool
	asset.Platform, ok = detector.DetectOS(fsConn)
	if !ok {
		return nil, errors.New("failed to detect OS")
	}
	asset.Id = conf.Type
	asset.Platform.Runtime = c.Runtime()
	return c, nil
}

func (c *AwsEbsConnection) FileInfo(path string) (shared.FileInfoDetails, error) {
	return shared.FileInfoDetails{}, errors.New("FileInfo not implemented")
}

func (c *AwsEbsConnection) FileSystem() afero.Fs {
	return c.FsProvider.FileSystem()
}

func (c *AwsEbsConnection) Close() {
	if c.opts != nil {
		if c.opts[snapshot.NoSetup] == "true" {
			return
		}
	}
	ctx := context.Background()
	err := c.volumeMounter.UnmountVolumeFromInstance()
	if err != nil {
		log.Error().Err(err).Msg("unable to unmount volume")
	}
	err = c.DetachVolumeFromInstance(ctx, c.scanVolumeInfo)
	if err != nil {
		log.Error().Err(err).Msg("unable to detach volume")
	}
	// only delete the volume if we created it, e.g., if we're scanning a snapshot
	if val, ok := c.scanVolumeInfo.Tags["createdBy"]; ok {
		if val == "Mondoo" {
			err = c.DeleteCreatedVolume(ctx, c.scanVolumeInfo)
			if err != nil {
				log.Error().Err(err).Msg("unable to delete volume")
			}
			log.Info().Str("vol-id", c.scanVolumeInfo.Id).Msg("deleted temporary volume created by Mondoo")
		}
	} else {
		log.Debug().Str("vol-id", c.scanVolumeInfo.Id).Msg("skipping volume deletion, not created by Mondoo")
	}
	err = c.volumeMounter.RemoveTempScanDir()
	if err != nil {
		log.Error().Err(err).Msg("unable to remove dir")
	}
}

func RawInstanceInfo(cfg aws.Config) (*imds.InstanceIdentityDocument, error) {
	metadata := imds.NewFromConfig(cfg)
	ctx := context.Background()
	doc, err := metadata.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		return nil, err
	}
	return &doc.InstanceIdentityDocument, nil
}

func (c *AwsEbsConnection) Identifier() (string, error) {
	return c.target.PlatformId, nil
}

func GetRawInstanceInfo(profile string) (*imds.InstanceIdentityDocument, error) {
	ctx := context.Background()
	var cfg aws.Config
	var err error
	if profile == "" {
		cfg, err = config.LoadDefaultConfig(ctx)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	}
	if err != nil {
		return nil, errors.Wrap(err, "could not load aws configuration")
	}
	i, err := RawInstanceInfo(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not load instance info: aws-ec2-ebs provider is only valid on ec2 instances")
	}
	return i, nil
}

func NewInstanceId(account string, region string, id string) (*awsec2ebstypes.InstanceId, error) {
	if region == "" || id == "" || account == "" {
		return nil, errors.New("invalid instance id. account, region and instance id required.")
	}
	return &awsec2ebstypes.InstanceId{Account: account, Region: region, Id: id}, nil
}

func ParseInstanceId(path string) (*awsec2ebstypes.InstanceId, error) {
	if !IsValidInstanceId(path) {
		return nil, errors.New("invalid instance id. expected account/<id>/region/<region-val>/instance/<instance-id>")
	}
	keyValues := strings.Split(path, "/")
	if len(keyValues) != 6 {
		return nil, errors.New("invalid instance id. expected account/<id>/region/<region-val>/instance/<instance-id>")
	}
	return NewInstanceId(keyValues[1], keyValues[3], keyValues[5])
}

var VALID_INSTANCE_ID = regexp.MustCompile(`^account/\d{12}/region\/(us(-gov)?|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\d\/instance\/.+$`)

func IsValidInstanceId(path string) bool {
	return VALID_INSTANCE_ID.MatchString(path)
}

func resourceTags(resourceType types.ResourceType, instanceId string) []types.TagSpecification {
	return []types.TagSpecification{
		{
			ResourceType: resourceType,
			Tags: []types.Tag{
				{Key: aws.String("createdBy"), Value: aws.String("Mondoo")},
				{Key: aws.String("Created By"), Value: aws.String("Mondoo")},
				{Key: aws.String("Created From Instance"), Value: aws.String(instanceId)},
			},
		},
	}
}

func ParseEbsTransportUrl(path string) (*awsec2ebstypes.EbsTransportTarget, error) {
	keyValues := strings.Split(path, "/")
	if len(keyValues) != 6 {
		return nil, errors.New("invalid id. expected account/<id>/region/<region-val>/{instance, volume, or snapshot}/<id>")
	}

	var itemType string
	switch keyValues[4] {
	case "volume":
		itemType = awsec2ebstypes.EBSTargetVolume
	case "snapshot":
		itemType = awsec2ebstypes.EBSTargetSnapshot
	default:
		itemType = awsec2ebstypes.EBSTargetInstance
	}

	return &awsec2ebstypes.EbsTransportTarget{Account: keyValues[1], Region: keyValues[3], Id: keyValues[5], Type: itemType}, nil
}

func (c *AwsEbsConnection) Name() string {
	return "aws ebs"
}

func (c *AwsEbsConnection) ID() uint32 {
	return c.id
}

func (c *AwsEbsConnection) Asset() *inventory.Asset {
	return c.asset
}

func (c *AwsEbsConnection) Capabilities() shared.Capabilities {
	return shared.Capability_RunCommand // not true, update to nothing
}

func (c *AwsEbsConnection) RunCommand(command string) (*shared.Command, error) {
	return nil, errors.New("unimplemented")
}

func (c *AwsEbsConnection) Type() shared.ConnectionType {
	return EBSConnectionType
}

func (c *AwsEbsConnection) Runtime() string {
	return "aws-ebs"
}

func (c *AwsEbsConnection) PlatformInfo() *inventory.Platform {
	return &inventory.Platform{
		Name:    "aws-ebs",
		Title:   "aws-ebs",
		Runtime: c.Runtime(),
	}
}
