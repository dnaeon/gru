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

// +build linux

package classifier

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
