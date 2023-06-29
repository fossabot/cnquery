package plugin

import (
	"github.com/hashicorp/go-plugin"
)

type Provider struct {
	Name       string
	Connectors []Connector
}

type Connector struct {
	Name      string
	Use       string   `json:",omitempty"`
	Short     string   `json:",omitempty"`
	Long      string   `json:",omitempty"`
	MinArgs   uint     `json:",omitempty"`
	MaxArgs   uint     `json:",omitempty"`
	Flags     []Flag   `json:",omitempty"`
	Aliases   []string `json:",omitempty"`
	Discovery []string `json:",omitempty"`
}

type FlagType byte

const (
	FlagType_Bool FlagType = 1 + iota
	FlagType_Int
	FlagType_String
	FlagType_List
	FlagType_KeyValue
)

type FlagOption byte

const (
	FlagOption_Hidden FlagOption = 0x1 << iota
	FlagOption_Deprecated
	FlagOption_Required
	FlagOption_Password
	// max: 8 options!
)

type Flag struct {
	Long    string     `json:",omitempty"`
	Short   string     `json:",omitempty"`
	Default string     `json:",omitempty"`
	Desc    string     `json:",omitempty"`
	Type    FlagType   `json:",omitempty"`
	Option  FlagOption `json:",omitempty"`
	// ConfigEntry that is used for this flag:
	// "" = use the same as Long
	// "some.other" = map to some.other field
	// "-" = do not read this from config
	ConfigEntry string `json:",omitempty"`
}

func Start(args []string, impl ProviderPlugin) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"provider": &ProviderPluginImpl{Impl: impl},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
