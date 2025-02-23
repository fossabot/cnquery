// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logindefs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mondoo.com/cnquery/providers/os/connection/mock"
	"go.mondoo.com/cnquery/providers/os/resources/logindefs"
)

func TestLoginDefsParser(t *testing.T) {
	mock, err := mock.New("./testdata/debian.toml", nil)
	require.NoError(t, err)

	f, err := mock.FileSystem().Open("/etc/login.defs")
	require.NoError(t, err)
	defer f.Close()

	entries := logindefs.Parse(f)

	assert.Equal(t, "tty", entries["TTYGROUP"])
	assert.Equal(t, "PATH=/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games", entries["ENV_PATH"])

	_, ok := entries["SHA_CRYPT_MIN_ROUNDS"]
	assert.False(t, ok)
}
