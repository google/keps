package hermetic

import (
	"fmt"
)

const (
	// a default email to use when committing on behalf of the user // TODO pass through principal user email
	kepToolEmail = "kubernetes-sig-contribex@googlegroups.com"

	// githubAuthorizationHeader is the header name for authorizing GitHub API requests
	githubAuthorizationHeaderName = "Authorization"
)

// githubGitUrl returns the URL to use for Git operations against a repo hosted on GitHub (e.g. clone, push)
func githubGitUrl(owner string, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

// githubForkUrl returns the URL to use when issuing a fork request to the GitHub API
func githubForkUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", owner, repo)
}

// githubPrUrl returns the URL to use when creating a pull request against a repo hosted on GitHub
func githubPrUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
}

// githubRepoApiUrl returns the URL to use when performing CRUD operations against a repo hosted on GitHub
func githubRepoApiUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
}
