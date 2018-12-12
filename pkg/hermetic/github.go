package hermetic

import (
	"fmt"
)

const (
	kepToolEmail                  = "kubernetes-sig-contribex@googlegroups.com"
	githubAuthorizationHeaderName = "Authorization"
)

func githubGitUrl(owner string, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
}

func githubForkUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", owner, repo)
}

func githubPrUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
}

func githubRepoApiUrl(owner string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
}
