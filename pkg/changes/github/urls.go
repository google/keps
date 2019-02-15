package github

import (
	"fmt"
)

// GitUrl returns the URL to use for Git operations against a repo hosted on GitHub (e.g. clone, pu
func GitUrl(owner string, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

// ForkUrl returns the URL to use when issuing a fork request to the GitHub API
func ForkUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", owner, repo)
}

// PrUrl returns the URL to use when creating a pull request against a repo hosted on GitHub
func PrUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
}

// RepoApiUrl returns the URL to use when performing CRUD operations against a repo hosted on GitHu
func RepoApiUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
}
