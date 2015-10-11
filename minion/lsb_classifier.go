// +build linux

package minion

import (
	"log"
	"bytes"
	"os/exec"
	"strings"
)

func init() {
	lsbdistid := NewCallbackClassifier("lsbdistid", "Distributor ID", lsbdistidClassifier)
	lsbdistdesc := NewCallbackClassifier("lsbdistdesc", "Short description of the distribution", lsbdistdescClassifier)
	lsbdistrelease := NewCallbackClassifier("lsbdistrelease", "Release number of the distribution", lsbdistreleaseClassifier)
	lsbdistcodename := NewCallbackClassifier("lsbdistcodename", "Distribution codename", lsbdistcodenameClassifier)

	RegisterClassifier(lsbdistid, lsbdistdesc, lsbdistrelease, lsbdistcodename)
}

func runLsbReleaseTool(args ...string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("/usr/bin/lsb_release", args...)
	cmd.Stdout = &buf
	err := cmd.Run()

	if err != nil {
		log.Printf("Failed to run lsb_release: %s\n", err)
		return "", err
	}

	data := strings.Split(buf.String(), ":")
	result := strings.TrimSpace(data[1])

	return result, nil
}

func lsbdistidClassifier() (string, error) {
	id, err := runLsbReleaseTool("--id")

	if err != nil {
		log.Println("Failed to get lsb dist id")
		return "", err
	}

	return id, nil
}

func lsbdistdescClassifier() (string, error) {
	desc, err := runLsbReleaseTool("--description")

	if err != nil {
		log.Println("Failed to get lsb dist description")
		return "", err
	}

	return desc, nil
}

func lsbdistreleaseClassifier() (string, error) {
	release, err := runLsbReleaseTool("--release")

	if err != nil {
		log.Println("Failed to get lsb dist release")
		return "", err
	}

	return release, nil
}

func lsbdistcodenameClassifier() (string, error) {
	codename, err := runLsbReleaseTool("--codename")

	if err != nil {
		log.Println("Failed to get lsb dist codename")
		return "", err
	}

	return codename, nil
}
