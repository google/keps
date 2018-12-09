package porcelain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	upstreamPRBaseName = "master" // TODO move into higher level function
)

type createPullRequestFunc func(string, string) (string, error)

func newCreatePRFunc(c *http.Client, token tokenProvider, apiUrl string, sourceLocation string) (createPullRequestFunc, error) {

	// we're not going to know when the pull request should actually be created
	// rather than store things like the API token directly on the Repository object
	// we create a callback designed to be integrated with a git operation (e.g push)
	// we're currently `leaking` a pull request because we don't create a handle to
	// remove it, but at the moment the only consumer would be a test
	var createPullRequestFunc = func(prTitle string, prDescription string) (string, error) {
		authToken, err := token.Value()
		if err != nil {
			return "", err
		}

		// serialize request payload
		var createPrPayload struct {
			TitleField          string `json:"title"`
			SourceBranchField   string `json:"head"`
			TargetBranchField   string `json:"base"`
			DescriptionField    string `json:"body,omitempty"`
			MaintainerCanModify bool   `json:"maintainer_can_modify,omitempty"`
		}

		createPrPayload.TitleField = prTitle
		createPrPayload.DescriptionField = prDescription
		createPrPayload.TargetBranchField = upstreamPRBaseName
		createPrPayload.SourceBranchField = sourceLocation

		payloadBytes, err := json.Marshal(createPrPayload)
		if err != nil {
			return "", err
		}

		payload := ioutil.NopCloser(bytes.NewBuffer(payloadBytes))

		// create HTTP request
		req, err := http.NewRequest(http.MethodPost, apiUrl, payload)
		if err != nil {
			return "", err
		}

		// add auth header
		req.Header.Add(githubAuthorizationHeaderName, fmt.Sprintf("token %s", authToken))

		// set context
		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		createPullRequest := req.WithContext(ctx)
		defer cancel()

		// Do request
		resp, err := c.Do(createPullRequest)
		if err != nil {
			return "", nil
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("expected status code 201 Created, got: %s", resp.Status)
		}

		// extract PR URL
		var createResponse struct {
			NodeIdField      string `json:"node_id"`
			PrNumberField    int    `json:"number"`
			HtmlUrlField     string `json:"html_url"`
			CommentsUrlField string `json:"comments_url"`
			PatchUrlField    string `json:"patch_url"`
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(bodyBytes, &createResponse)
		if err != nil {
			return "", err
		}

		// return PR URL
		return createResponse.HtmlUrlField, nil
	}

	return createPullRequestFunc, nil
}
