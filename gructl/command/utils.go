// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package command

import (
	"regexp"
	"strings"

	"github.com/dnaeon/gru/client"

	etcdclient "github.com/coreos/etcd/client"
	"github.com/pborman/uuid"
	"github.com/urfave/cli"
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
