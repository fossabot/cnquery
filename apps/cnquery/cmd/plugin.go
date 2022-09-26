package cmd

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-plugin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.mondoo.com/cnquery/cli/assetlist"
	"go.mondoo.com/cnquery/cli/printer"
	"go.mondoo.com/cnquery/cli/shell"
	"go.mondoo.com/cnquery/cli/theme"
	"go.mondoo.com/cnquery/logger"
	"go.mondoo.com/cnquery/motor/asset"
	"go.mondoo.com/cnquery/motor/discovery"
	"go.mondoo.com/cnquery/motor/inventory"
	provider_resolver "go.mondoo.com/cnquery/motor/providers/resolver"
	"go.mondoo.com/cnquery/mqlc"
	"go.mondoo.com/cnquery/mqlc/parser"
	"go.mondoo.com/cnquery/resources/packs/os/info"
	"go.mondoo.com/cnquery/shared"
	"go.mondoo.com/cnquery/shared/proto"
)

// pluginCmd represents the version command
var pluginCmd = &cobra.Command{
	Use:    "run_as_plugin",
	Hidden: true,
	Short:  "Run as a plugin",
	Run: func(cmd *cobra.Command, args []string) {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: shared.Handshake,
			Plugins: map[string]plugin.Plugin{
				"counter": &shared.CNQueryPlugin{Impl: &cnqueryPlugin{}},
			},

			// A non-nil value here enables gRPC serving for this plugin...
			GRPCServer: plugin.DefaultGRPCServer,
		})
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}

type cnqueryPlugin struct{}

func (c *cnqueryPlugin) RunQuery(conf *proto.RunQueryConfig, out shared.OutputHelper) error {
	if conf.Command == "" {
		return errors.New("No command provided, nothing to do.")
	}

	ctx := discovery.InitCtx(context.Background())

	if conf.DoParse {
		ast, err := parser.Parse(conf.Command)
		if err != nil {
			return errors.Wrap(err, "failed to parse command")
		}
		out.WriteString(logger.PrettyJSON(ast))
		return nil
	}

	if conf.DoAst {
		b, err := mqlc.Compile(conf.Command, info.Registry.Schema(), conf.Features, nil)
		if err != nil {
			return errors.Wrap(err, "failed to compile command")
		}

		out.WriteString(logger.PrettyJSON((b)) + "\n" + printer.DefaultPrinter.CodeBundle(b))

		return nil
	}

	log.Info().Msgf("discover related assets for %d asset(s)", len(conf.Inventory.Spec.Assets))
	im, err := inventory.New(inventory.WithInventory(conf.Inventory))
	if err != nil {
		return errors.Wrap(err, "could not load asset information")
	}
	assetErrors := im.Resolve(ctx)
	if len(assetErrors) > 0 {
		for a := range assetErrors {
			log.Error().Err(assetErrors[a]).Str("asset", a.Name).Msg("could not resolve asset")
		}
	}

	assetList := im.GetAssets()
	if len(assetList) == 0 {
		return errors.New("could not find an asset that we can connect to")
	}

	var connectAsset *asset.Asset

	if len(assetList) == 1 {
		connectAsset = assetList[0]
	} else if len(assetList) > 1 && conf.PlatformId != "" {
		connectAsset, err = filterAssetByPlatformID(assetList, conf.PlatformId)
		if err != nil {
			return err
		}
	} else if len(assetList) > 1 {
		r := assetlist.NewSimpleRenderer(theme.OperatingSytemTheme)
		out.WriteString(r.Render(assetList) + "\n")
		return errors.New("cannot connect to more than one asset, use --platform-id to select a specific asset")
	}

	if conf.DoRecord {
		log.Info().Msg("enable recording of platform calls")
	}

	m, err := provider_resolver.OpenAssetConnection(ctx, connectAsset, im.GetCredential, conf.DoRecord)
	if err != nil {
		return errors.New("could not connect to asset")
	}

	// when we close the shell, we need to close the backend and store the recording
	onCloseHandler := func() {
		storeRecording(m)
	}

	shellOptions := []shell.ShellOption{}
	shellOptions = append(shellOptions, shell.WithOnCloseListener(onCloseHandler))
	shellOptions = append(shellOptions, shell.WithFeatures(conf.Features))
	shellOptions = append(shellOptions, shell.WithOutput(out))

	sh, err := shell.New(m, shellOptions...)
	if err != nil {
		return errors.Wrap(err, "failed to initialize the shell")
	}
	defer sh.Close()

	code, results, err := sh.RunOnce(conf.Command)
	if err != nil {
		return errors.Wrap(err, "failed to run")
	}

	if conf.Format != "json" {
		sh.PrintResults(code, results)
		return nil
	}

	var checksums []string
	eps := code.CodeV2.Entrypoints()
	checksums = make([]string, len(eps))
	for i, ref := range eps {
		checksums[i] = code.CodeV2.Checksums[ref]
	}

	for _, checksum := range checksums {
		result := results[checksum]
		if result == nil {
			return errors.New("cannot find result for this query")
		}

		if result.Data.Error != nil {
			return result.Data.Error
		}

		j := result.Data.JSON(checksum, code)
		out.Write(append(j, '\n'))
	}

	return nil
}
