package porcelain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type createForkRequestFunc func() (string, string, error)

func newCreateForkFunc(c *http.Client, token tokenProvider, apiUrl string) (createForkRequestFunc, error) {

	var forkRepoFunc = func() (string, string, error) {
		authToken, err := token.Value()
		if err != nil {
			return "", "", err
		}

		req, err := http.NewRequest(http.MethodPost, apiUrl, nil)
		if err != nil {
			return "", "", err
		}

		// add auth header
		req.Header.Add(githubAuthorizationHeaderName, fmt.Sprintf("token %s", authToken))

		// set context
		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		forkRepo := req.WithContext(ctx)
		defer cancel()

		// Do request
		resp, err := c.Do(forkRepo)
		if err != nil {
			return "", "", nil
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted {
			return "", "", fmt.Errorf("expected status code 202 Accepted, got: %s", resp.Status)
		}

		// extract API url
		var forkResponse struct {
			NodeIdField  string `json:"node_id"`
			HtmlUrlField string `json:"html_url"`
			ApiUrlField  string `json:"url"`
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}

		err = json.Unmarshal(bodyBytes, &forkResponse)
		if err != nil {
			return "", "", err
		}

		return forkResponse.ApiUrlField, forkResponse.HtmlUrlField, nil
	}

	return forkRepoFunc, nil
}
