// +build linux

package minion

import (
	"bytes"
	"os/exec"
	"strings"
)

func init() {
	Register("lsbdistid", idProvider)
	Register("lsbdistdesc", descProvider)
	Register("lsbdistrelease", releaseProvider)
	Register("lsbdistcodename", codenameProvider)
}

func runLSBreleaseTool(args ...string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("/usr/bin/lsb_release", args...)
	cmd.Stdout = &buf
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	data := strings.Split(buf.String(), ":")
	result := strings.TrimSpace(data[1])

	return result, nil
}

func idProvider() (string, error) {
	return runLSBreleaseTool("--id")
}

func descProvider() (string, error) {
	return runLSBreleaseTool("--description")
}

func releaseProvider() (string, error) {
	return runLSBreleaseTool("--release")
}

func codenameProvider() (string, error) {
	return runLSBreleaseTool("--codename")
}
