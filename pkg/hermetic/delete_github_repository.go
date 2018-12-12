package hermetic

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type deleteUserRepoFunc func() error

func newDeleteGithubUserRepoFunc(c *http.Client, token tokenProvider, apiUrl string) (deleteUserRepoFunc, error) {
	var deleteRepo = func() error {
		authToken, err := token.Value()
		if err != nil {
			return err
		}

		deleteReq, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
		if err != nil {
			return err
		}

		deleteReq.Header.Add(githubAuthorizationHeaderName, fmt.Sprintf("token %s", authToken))

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

	return deleteRepo, nil
}
