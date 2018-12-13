package hermetic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type createForkRequestFunc func() error

// newCreateForkFunc creates a callback that is able to handle the details of issuing a GitHub API request
// to fork a repository into the authenticated user's GitHub account (forking into an organization is not
// supported at this time
func newCreateForkFunc(c *http.Client, token tokenProvider, apiUrl string) (createForkRequestFunc, error) {

	var forkRepoFunc = func() error {
		authToken, err := token.Value()
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, apiUrl, nil)
		if err != nil {
			return err
		}

		req.Header.Add(githubAuthorizationHeaderName, fmt.Sprintf("token %s", authToken))

		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		forkRepo := req.WithContext(ctx)
		defer cancel()

		resp, err := c.Do(forkRepo)
		if err != nil {
			return nil
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("expected status code 202 Accepted, got: %s", resp.Status)
		}

		var forkResponse struct {
			NodeIdField  string `json:"node_id"`
			HtmlUrlField string `json:"html_url"`
			ApiUrlField  string `json:"url"`
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bodyBytes, &forkResponse)
		if err != nil {
			return err
		}

		if forkResponse.ApiUrlField == "" {
			return fmt.Errorf("recieved empty API url from response: \n %s", string(bodyBytes))
		}

		if forkResponse.HtmlUrlField == "" {
			return fmt.Errorf("recieved empty Git url from response: \n %s", string(bodyBytes))
		}

		return nil
	}

	return forkRepoFunc, nil
}
