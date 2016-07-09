package utils

import "os/exec"

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
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	repo := &GitRepo{
		Path:     path,
		Upstream: upstream,
		git:      path,
	}

	return repo, nil
}

// Fetch fetches from the given remote
func (gr *GitRepo) Fetch(remote string) ([]byte, error) {
	return exec.Command(gr.git, "--git-dir", gr.Path, "fetch", remote).CombinedOutput()
}

// Pull pulls from the given remote and merges changes into the
// local branch
func (gr *GitRepo) Pull(remote, branch string) ([]byte, error) {
	out, err := gr.Checkout(branch)
	if err != nil {
		return out, err
	}

	return exec.Command(gr.git, "--git-dir", gr.Path, "pull", remote).CombinedOutput()
}

// Checkout checks out a given local branch
func (gr *GitRepo) Checkout(branch string) ([]byte, error) {
	return exec.Command(gr.git, "--git-dir", gr.Path, "checkout", branch).CombinedOutput()
}

// CheckoutDetached checks out a given local branch in detached mode
func (gr *GitRepo) CheckoutDetached(branch string) ([]byte, error) {
	return exec.Command(gr.git, "--git-dir", gr.Path, "checkout", "--detach", branch).CombinedOutput()
}

// Clone clones the upstream repository
func (gr *GitRepo) Clone() ([]byte, error) {
	return exec.Command(gr.git, "clone", gr.Upstream, gr.Path).CombinedOutput()
}

// Head returns the SHA1 commit id at HEAD
func (gr *GitRepo) Head() ([]byte, error) {
	return exec.Command(gr.git, "--git-dir", gr.Path, "rev-parse", "--short", "HEAD").CombinedOutput()
}

// IsGitRepo checks if the repository is a valid Git repository
func (gr *GitRepo) IsGitRepo() bool {
	err := exec.Command(gr.git, "--git-dir", gr.Path, "rev-parse").Run()
	if err != nil {
		return false
	}

	return true
}
