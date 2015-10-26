package command

import (
	"strings"

	"github.com/dnaeon/gru/client"

	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
)

func newEtcdMinionClientFromFlags(c *cli.Context) client.Client {
	eFlag := c.GlobalString("endpoint")
	uFlag := c.GlobalString("username")
	pFlag := c.GlobalString("password")
	tFlag := c.GlobalDuration("timeout")

	cfg := etcdclient.Config{
		Endpoints: strings.Split(eFlag, ","),
		Transport: etcdclient.DefaultTransport,
		HeaderTimeoutPerRequest: tFlag,
	}

	if uFlag != "" && pFlag != "" {
		cfg.Username = uFlag
		cfg.Password = pFlag
	}

	klient := client.NewEtcdMinionClient(cfg)

	return klient
}
