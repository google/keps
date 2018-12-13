package hermetic

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type pushRepoChangesFunc func() error

// newPushRepoChanges creates a callback to abstract away the details of pushing local changes to a remote Git repository
// named "origin"
func newPushRepoChangesFunc(gitRepo *git.Repository, token tokenProvider, withBranchName string) (pushRepoChangesFunc, error) {

	var pushRepoChangesFunc = func() error {

		authToken, err := token.Value()
		if err != nil {
			return err
		}

		worktree, err := gitRepo.Worktree()
		if err != nil {
			return err
		}

		repoStatus, err := worktree.Status()
		if err != nil {
			return err
		}

		uncommittedFiles := !repoStatus.IsClean()
		if uncommittedFiles {
			return fmt.Errorf("uncommitted files exist. Please commit or remove them before pushing changes: %s", repoStatus)
		}

		// push local changes
		err = gitRepo.Push(&git.PushOptions{
			RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", withBranchName, withBranchName))},
			RemoteName: OriginRemoteName,
			Auth: &http.BasicAuth{
				Username: arbitraryBasicAuthUsername,
				Password: authToken,
			},
		})

		switch err {
		case nil:
			return nil
		case git.NoErrAlreadyUpToDate:
			return nil
		default:
			return err
		}

	}

	return pushRepoChangesFunc, nil
}
