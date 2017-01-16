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

package utils

import (
	"os/exec"
	"strings"
)

// GitRepo type manages a VCS repository with Git
type GitRepo struct {
	// Local path to the repository
	Path string

	// Upstream URL of the Git repository
	Upstream string

	// Path to the Git tool
	git string
}

// NewGitRepo creates a new Git repository
func NewGitRepo(path, upstream string) (*GitRepo, error) {
	cmd, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	repo := &GitRepo{
		Path:     path,
		Upstream: upstream,
		git:      cmd,
	}

	return repo, nil
}

// Fetch fetches from the given remote
func (gr *GitRepo) Fetch(remote string) ([]byte, error) {
	return exec.Command(gr.git, "-C", gr.Path, "fetch", remote).CombinedOutput()
}

// Pull pulls from the given remote and merges changes into the
// local branch
func (gr *GitRepo) Pull(remote, branch string) ([]byte, error) {
	out, err := gr.Checkout(branch)
	if err != nil {
		return out, err
	}

	return exec.Command(gr.git, "-C", gr.Path, "pull", remote).CombinedOutput()
}

// Checkout checks out a given local branch
func (gr *GitRepo) Checkout(branch string) ([]byte, error) {
	return exec.Command(gr.git, "-C", gr.Path, "checkout", branch).CombinedOutput()
}

// CheckoutDetached checks out a given local branch in detached mode
func (gr *GitRepo) CheckoutDetached(branch string) ([]byte, error) {
	return exec.Command(gr.git, "-C", gr.Path, "checkout", "--detach", branch).CombinedOutput()
}

// Clone clones the upstream repository
func (gr *GitRepo) Clone() ([]byte, error) {
	return exec.Command(gr.git, "clone", gr.Upstream, gr.Path).CombinedOutput()
}

// Head returns the SHA1 commit id at HEAD
func (gr *GitRepo) Head() (string, error) {
	head, err := exec.Command(gr.git, "-C", gr.Path, "rev-parse", "--short", "HEAD").CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(head), "\n"), nil
}

// IsGitRepo checks if the repository is a valid Git repository
func (gr *GitRepo) IsGitRepo() bool {
	err := exec.Command(gr.git, "-C", gr.Path, "rev-parse").Run()
	if err != nil {
		return false
	}

	return true
}
