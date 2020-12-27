package cmd

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"

	"github.com/spf13/cobra"

	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

const (
	icaImport = "github.com/chainapsis/cosmos-sdk-interchain-account"
	apppkg    = "app"
	// moduleDir        = "x"
	// icaVersionCommit = "daba1321259a442f929f82738b9a9d632eeb4351"

	// Placeholders in Stargate app.go
	placeholderSgAppModuleImport      = "// this line is used by starport scaffolding # stargate/app/moduleImport"
	placeholderSgAppModuleBasic       = "// this line is used by starport scaffolding # stargate/app/moduleBasic"
	placeholderSgAppKeeperDeclaration = "// this line is used by starport scaffolding # stargate/app/keeperDeclaration"
	placeholderSgAppStoreKey          = "// this line is used by starport scaffolding # stargate/app/storeKey"
	placeholderSgAppKeeperDefinition  = "// this line is used by starport scaffolding # stargate/app/keeperDefinition"
	placeholderSgAppAppModule         = "// this line is used by starport scaffolding # stargate/app/appModule"
	placeholderSgAppInitGenesis       = "// this line is used by starport scaffolding # stargate/app/initGenesis"
	placeholderSgAppParamSubspace     = "// this line is used by starport scaffolding # stargate/app/paramSubspace"
)

func NewModule() *cobra.Command {
	c := &cobra.Command{
		Use:   "module",
		Short: "Manage ICA module for cosmos app",
	}
	c.AddCommand(NewModuleImport())
	return c
}

func NewModuleImport() *cobra.Command {
	c := &cobra.Command{
		Use:   "import",
		Short: "Import a ICA modulr to cosmos app",
		RunE:  importModuleHandler,
	}
	return c
}

func importModuleHandler(cmd *cobra.Command, args []string) error {
	version, err := cosmosver.Detect("")
	if err != nil {
		return err
	}

	if version != cosmosver.Stargate {
		return fmt.Errorf("ICA only support the stargate")
	}

	installed, err := isICAImported("")
	if err != nil {
		return err
	}
	if installed && len(args) == 0 {
		return nil
	}

	if !installed {
		err = installICA()
		if err != nil {
			return err
		}
	}

	if len(args) > 0 {
		if args[0] == "mock" {
			return importMockModule()
		}
		return fmt.Errorf("unknown module")
	}

	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		path := "app/app.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	ibcaccount "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account"
	ibcaccountkeeper "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	ibcaccounttypes "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"`
		replacement := fmt.Sprintf(template, placeholderSgAppModuleImport)
		content := strings.Replace(f.String(), placeholderSgAppModuleImport, replacement, 1)

		template2 := `%[1]v
		ibcaccount.AppModuleBasic{},`
		replacement2 := fmt.Sprintf(template2, placeholderSgAppModuleBasic)
		content = strings.Replace(content, placeholderSgAppModuleBasic, replacement2, 1)

		template3 := `%[1]v
	IBCAccountKeeper     ibcaccountkeeper.Keeper
	ScopedIBCAccountKeeper capabilitykeeper.ScopedKeeper`
		replacement3 := fmt.Sprintf(template3, placeholderSgAppKeeperDeclaration)
		content = strings.Replace(content, placeholderSgAppKeeperDeclaration, replacement3, 1)

		template4 := `%[1]v
		ibcaccounttypes.StoreKey,`
		replacement5 := fmt.Sprintf(template4, placeholderSgAppStoreKey)
		content = strings.Replace(content, placeholderSgAppStoreKey, replacement5, 1)

		template5 := `scopedIBCAccountKeeper := app.CapabilityKeeper.ScopeToModule(ibcaccounttypes.ModuleName)
	app.IBCAccountKeeper = ibcaccountkeeper.NewKeeper(appCodec, keys[ibcaccounttypes.StoreKey],
	map[string]ibcaccounttypes.TxEncoder{
		// register the tx encoder for cosmos-sdk
		"cosmos-sdk": ibcaccountkeeper.SerializeCosmosTx(appCodec, interfaceRegistry),
	}, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
	app.AccountKeeper, scopedIBCAccountKeeper, app.Router(),
	)
	ibcAccountModule := ibcaccount.NewAppModule(app.IBCAccountKeeper)
	ibcRouter.AddRoute(ibcaccounttypes.ModuleName, ibcAccountModule)
	app.IBCKeeper.SetRouter(ibcRouter)`
		content = strings.Replace(content, "app.IBCKeeper.SetRouter(ibcRouter)", template5, 1)

		template6 := `%[1]v
		ibcAccountModule,`
		replacement6 := fmt.Sprintf(template6, placeholderSgAppAppModule)
		content = strings.Replace(content, placeholderSgAppAppModule, replacement6, 1)

		template7 := `%[1]v
		ibcaccounttypes.ModuleName,`
		replacement7 := fmt.Sprintf(template7, placeholderSgAppInitGenesis)
		content = strings.Replace(content, placeholderSgAppInitGenesis, replacement7, 1)

		template8 := `%[1]v
	paramsKeeper.Subspace(ibcaccounttypes.ModuleName)`
		replacement8 := fmt.Sprintf(template8, placeholderSgAppParamSubspace)
		content = strings.Replace(content, placeholderSgAppParamSubspace, replacement8, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	})

	run := genny.WetRunner(context.Background())
	run.With(g)

	return run.Run()
}

func importMockModule() error {
	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		path := "app/app.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	ibcaccountmock "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock"
	ibcaccountmockkeeper "github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/keeper"`
		replacement := fmt.Sprintf(template, placeholderSgAppModuleImport)
		content := strings.Replace(f.String(), placeholderSgAppModuleImport, replacement, 1)

		template2 := `%[1]v
		ibcaccountmock.AppModuleBasic{},`
		replacement2 := fmt.Sprintf(template2, placeholderSgAppModuleBasic)
		content = strings.Replace(content, placeholderSgAppModuleBasic, replacement2, 1)

		template3 := `%[1]v
	IBCAccountMockKeeper ibcaccountmockkeeper.Keeper`
		replacement3 := fmt.Sprintf(template3, placeholderSgAppKeeperDeclaration)
		content = strings.Replace(content, placeholderSgAppKeeperDeclaration, replacement3, 1)

		template5 := `%[1]v
	app.IBCAccountMockKeeper = ibcaccountmockkeeper.NewKeeper(app.IBCAccountKeeper)
	ibcAccountMockModule := ibcaccountmock.NewAppModule(app.IBCAccountMockKeeper)`
		replacement5 := fmt.Sprintf(template5, placeholderSgAppKeeperDefinition)
		content = strings.Replace(content, placeholderSgAppKeeperDefinition, replacement5, 1)

		template6 := `%[1]v
		ibcAccountMockModule,`
		replacement6 := fmt.Sprintf(template6, placeholderSgAppAppModule)
		content = strings.Replace(content, placeholderSgAppAppModule, replacement6, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	})

	run := genny.WetRunner(context.Background())
	run.With(g)

	return run.Run()
}

func isICAImported(appPath string) (bool, error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, apppkg))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ImportsOnly)
	if err != nil {
		return false, err
	}
	for _, pkg := range all {
		for _, f := range pkg.Files {
			for _, imp := range f.Imports {
				if strings.Contains(imp.Path.Value, icaImport) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func installICA() error {
	return cmdrunner.
		New(
			cmdrunner.DefaultStderr(os.Stderr),
		).
		Run(context.Background(),
			step.New(
				step.Exec(
					"go",
					"get",
					icaImport,
				),
			),
		)
}
