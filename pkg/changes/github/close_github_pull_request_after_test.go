package github_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/calebamiles/keps/pkg/changes/auth"
	"github.com/calebamiles/keps/pkg/changes/github"
)

func closeGithubPullRequest(token auth.TokenProvider, pullRequestURL string) error {
	c := github.HttpClient()

	apiUrl := strings.Replace(pullRequestURL, "github.com", "api.github.com/repos", 1)
	apiUrl = strings.Replace(apiUrl, "pull", "pulls", 1)

	var payloadBody struct {
		State string `json:"state"`
	}

	payloadBody.State = "closed"
	payloadBytes, err := json.Marshal(payloadBody)
	if err != nil {
		return err
	}

	payload := ioutil.NopCloser(bytes.NewBuffer(payloadBytes))

	// create HTTP request
	req, err := http.NewRequest(http.MethodPatch, apiUrl, payload)
	if err != nil {
		return err
	}

	// add auth header
	err = github.AddAuthorizationHeader(req, token)
	if err != nil {
		return err
	}

	// set context
	ctx, cancel := context.WithTimeout(context.Background(), 7000*time.Millisecond)
	closePullRequest := req.WithContext(ctx)
	defer cancel()

	// Do request
	resp, err := c.Do(closePullRequest)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("when closing pull request, expected 200 Status OK, got: %s", resp.Status)
	}

	return nil
}
