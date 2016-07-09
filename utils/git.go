package utils

// GitRepo type manages a VCS repository with Git
type GitRepo struct {
	// Local path to the repository
	Path string

	// Upstream URL of the Git repository
	Upstream string
}

// NewGitRepo creates a new Git repository
func NewGitRepo(path, upstream string) *GitRepo {
	return &GitRepo{
		Path:     path,
		Upstream: upstream,
	}
}
