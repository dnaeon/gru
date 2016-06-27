package command

import (
	"github.com/dnaeon/gru/catalog"
	"github.com/dnaeon/gru/resource"
	"github.com/urfave/cli"
	"github.com/yuin/gopher-lua"
)

// NewApplyCommand creates a new sub-command for
// applying configurations on the local system
func NewApplyCommand() cli.Command {
	cmd := cli.Command{
		Name:   "apply",
		Usage:  "apply configuration on the local system",
		Action: execApplyCommand,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "siterepo",
				Value:  "",
				Usage:  "path/url to the site repo",
				EnvVar: "GRU_SITEREPO",
			},
			cli.BoolFlag{
				Name:  "dry-run",
				Usage: "just report what would be done, instead of doing it",
			},
		},
	}

	return cmd
}

// Executes the "apply" command
func execApplyCommand(c *cli.Context) error {
	if len(c.Args()) < 1 {
		return cli.NewExitError(errNoModuleName.Error(), 64)
	}

	L := lua.NewState()
	defer L.Close()
	config := &catalog.Config{
		Module:   c.Args()[0],
		DryRun:   c.Bool("dry-run"),
		Logger:   resource.DefaultLogger,
		SiteRepo: c.String("siterepo"),
		L:        L,
	}

	katalog := catalog.New(config)
	if err := katalog.Load(); err != nil {
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	if err := katalog.Run(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
