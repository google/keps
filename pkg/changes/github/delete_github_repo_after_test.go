package github_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/calebamiles/keps/pkg/changes/auth"
	"github.com/calebamiles/keps/pkg/changes/github"
)

// one day we may have an actual production use for this functionality but until then it'll live here
func deleteGithubRepo(token auth.TokenProvider, owner string, repo string) error {
	c := github.HttpClient()
	apiUrl := github.RepoApiUrl(owner, repo)

	deleteReq, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
	if err != nil {
		return err
	}

	err = github.AddAuthorizationHeader(deleteReq, token)
	if err != nil {
		return err
	}

	delCtx, cancelDel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	defer cancelDel()

	deleteRepo := deleteReq.WithContext(delCtx)

	delResp, delErr := c.Do(deleteRepo)
	if delErr != nil {
		return delErr
	}

	delErr = delResp.Body.Close()
	if delErr != nil {
		return delErr
	}

	if delResp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected Status: 204 No Content when deleting repository; got: %s. URL: %s", delResp.Status, apiUrl)
	}

	return nil
}
