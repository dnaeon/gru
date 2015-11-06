package classifier

import "os"

func init() {
	Register("fqdn", fqdnProvider)
}

func fqdnProvider() (string, error) {
	hostname, err := os.Hostname()

	return hostname, err
}
