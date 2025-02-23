// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mondoo.com/cnquery/providers/os/connection/mock"
	"go.mondoo.com/cnquery/providers/os/detector"
)

func TestDetectLinuxInstance(t *testing.T) {
	conn, err := mock.New("./testdata/instance_linux.toml", nil)
	require.NoError(t, err)
	platform, ok := detector.DetectOS(conn)
	require.True(t, ok)

	identifier, name, relatedIdentifiers := Detect(conn, platform)

	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/compute/v1/projects/mondoo-dev-262313/zones/us-central1-a/instances/6001244637815193808", identifier)
	assert.Equal(t, "", name)
	require.Len(t, relatedIdentifiers, 1)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/projects/mondoo-dev-262313", relatedIdentifiers[0])
}

func TestDetectWindowsInstance(t *testing.T) {
	conn, err := mock.New("./testdata/instance_windows.toml", nil)
	require.NoError(t, err)
	platform, ok := detector.DetectOS(conn)
	require.True(t, ok)

	identifier, name, relatedIdentifiers := Detect(conn, platform)

	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/compute/v1/projects/mondoo-dev-262313/zones/us-central1-a/instances/5275377306317132843", identifier)
	assert.Equal(t, "", name)
	require.Len(t, relatedIdentifiers, 1)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/projects/mondoo-dev-262313", relatedIdentifiers[0])
}

func TestNoMatch(t *testing.T) {
	conn, err := mock.New("./testdata/aws_instance.toml", nil)
	require.NoError(t, err)
	platform, ok := detector.DetectOS(conn)
	require.True(t, ok)

	identifier, name, relatedIdentifiers := Detect(conn, platform)

	assert.Empty(t, identifier)
	assert.Empty(t, name)
	assert.Empty(t, relatedIdentifiers)
}
