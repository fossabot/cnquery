package platformid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mondoo.io/mondoo/motor/providers/mock"
)

func TestGuidWindows(t *testing.T) {
	provider, err := mock.NewFromTomlFile("./testdata/guid_windows.toml")
	require.NoError(t, err)

	lid := WinIdProvider{provider: provider}
	id, err := lid.ID()
	require.NoError(t, err)

	assert.Equal(t, "6BAB78BE-4623-4705-924C-2B22433A4489", id)
}
