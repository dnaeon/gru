package command

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dnaeon/gru/minion"
	"github.com/urfave/cli"
)

// NewServeCommand creates a new sub-command for starting a
// minion and its services
func NewServeCommand() cli.Command {
	cmd := cli.Command{
		Name:   "serve",
		Usage:  "start minion",
		Action: execServeCommand,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "concurrency",
				Usage: "number of goroutines used for concurrent processing",
				Value: runtime.NumCPU(),
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "set minion name",
				Value: "",
			},
			cli.StringFlag{
				Name:   "siterepo",
				Value:  "",
				Usage:  "path/url to the site repo",
				EnvVar: "GRU_SITEREPO",
			},
		},
	}

	return cmd
}

// Executes the "serve" command
func execServeCommand(c *cli.Context) error {
	name, err := os.Hostname()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	concurrency := c.Int("concurrency")
	if concurrency < 0 {
		concurrency = runtime.NumCPU()
	}

	if c.String("siterepo") == "" {
		return cli.NewExitError(errNoSiteRepo.Error(), 64)
	}

	nameFlag := c.String("name")
	if nameFlag != "" {
		name = nameFlag
	}

	etcdCfg := etcdConfigFromFlags(c)
	minionCfg := &minion.EtcdMinionConfig{
		Concurrency: concurrency,
		Name:        name,
		SiteRepo:    c.String("siterepo"),
		EtcdConfig:  etcdCfg,
	}

	m, err := minion.NewEtcdMinion(minionCfg)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// Channel on which the shutdown signal is sent
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start minion
	err = m.Serve()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// Block until a shutdown signal is received
	<-quit
	m.Stop()

	return nil
}
