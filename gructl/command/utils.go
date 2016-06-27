package command

import (
	"regexp"
	"strings"

	"github.com/dnaeon/gru/client"

	"github.com/urfave/cli"
	etcdclient "github.com/coreos/etcd/client"
	"github.com/pborman/uuid"
)

func etcdConfigFromFlags(c *cli.Context) etcdclient.Config {
	eFlag := c.GlobalString("endpoint")
	uFlag := c.GlobalString("username")
	pFlag := c.GlobalString("password")
	tFlag := c.GlobalDuration("timeout")

	cfg := etcdclient.Config{
		Endpoints:               strings.Split(eFlag, ","),
		Transport:               etcdclient.DefaultTransport,
		HeaderTimeoutPerRequest: tFlag,
	}

	if uFlag != "" && pFlag != "" {
		cfg.Username = uFlag
		cfg.Password = pFlag
	}

	return cfg
}

func newEtcdMinionClientFromFlags(c *cli.Context) client.Client {
	cfg := etcdConfigFromFlags(c)
	klient := client.NewEtcdMinionClient(cfg)

	return klient
}

// Parses a classifier pattern and returns
// minions which match the given classifier pattern.
// A classifier pattern is described as 'key=regexp',
// where 'key' is a classifier key and 'regexp' is a
// regular expression that is compiled and matched
// against the minions' classifier values.
// If 'key' is empty all registered minions are returned.
// If 'regexp' is empty or missing all minions which
// contain the given 'key' are returned instead.
func parseClassifierPattern(klient client.Client, pattern string) ([]uuid.UUID, error) {
	// If no classifier pattern provided,
	// return all registered minions
	if pattern == "" {
		return klient.MinionList()
	}

	data := strings.SplitN(pattern, "=", 2)
	key := data[0]

	// If only a classifier key is provided, return all
	// minions which contain the given classifier key
	if len(data) == 1 {
		return klient.MinionWithClassifierKey(key)
	}

	toMatch := data[1]
	re, err := regexp.Compile(toMatch)
	if err != nil {
		return nil, err
	}

	minions, err := klient.MinionWithClassifierKey(key)
	if err != nil {
		return nil, err
	}

	var result []uuid.UUID
	for _, minion := range minions {
		klassifier, err := klient.MinionClassifier(minion, key)
		if err != nil {
			return nil, err
		}

		if re.MatchString(klassifier.Value) {
			result = append(result, minion)
		}
	}

	return result, nil
}
