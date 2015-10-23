package command

import (
	"strings"

	"github.com/dnaeon/gru/client"

	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
)

func newEtcdMinionClientFromFlags(c *cli.Context) client.Client {
	endpointFlag := c.GlobalString("endpoint")
	usernameFlag := c.GlobalString("username")
	passwordFlag := c.GlobalString("password")
	timeoutFlag := c.GlobalDuration("timeout")

	cfg := etcdclient.Config{
		Endpoints: strings.Split(endpointFlag, ","),
		Transport: etcdclient.DefaultTransport,
		HeaderTimeoutPerRequest: timeoutFlag,
	}

	if usernameFlag != "" && passwordFlag != "" {
		cfg.Username = usernameFlag
		cfg.Password = passwordFlag
	}

	klient := client.NewEtcdMinionClient(cfg)

	return klient
}
