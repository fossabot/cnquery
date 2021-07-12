package transports

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestCredentialParser(t *testing.T) {
	content := `
- type: password
  user: username
  password: pass
- type: private_key
  user: username
  identity_file: /path/to/key
  password: password
- type: credentials_query
`

	v := []*Credential{}
	yaml.Unmarshal([]byte(content), &v)

	assert.Equal(t, 3, len(v))
	assert.Equal(t, CredentialType_password, v[0].Type)
	assert.Equal(t, CredentialType_private_key, v[1].Type)
	assert.Equal(t, CredentialType_credentials_query, v[2].Type)
}

func TestCredentialMarshal(t *testing.T) {
	data, err := json.Marshal(CredentialType_undefined)
	require.NoError(t, err)
	assert.Equal(t, "\"undefined\"", string(data))
}
