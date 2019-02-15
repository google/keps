package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/calebamiles/keps/pkg/changes/auth"
)

// ForkRepository forks a GitHub repository, returning a clone-able git URL
func Fork(token auth.TokenProvider, owner string, repo string) (string, error) {
	apiUrl := ForkUrl(owner, repo)

	req, err := http.NewRequest(http.MethodPost, apiUrl, nil)
	if err != nil {
		return "", err
	}

	err = AddAuthorizationHeader(req, token)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	forkRepo := req.WithContext(ctx)
	defer cancel()

	c := HttpClient()
	resp, err := c.Do(forkRepo)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("expected status code 202 Accepted, got: %s", resp.Status)
	}

	var forkResponse struct {
		NodeIdField  string `json:"node_id"`
		HtmlUrlField string `json:"html_url"`
		ApiUrlField  string `json:"url"`
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bodyBytes, &forkResponse)
	if err != nil {
		return "", err
	}

	if forkResponse.ApiUrlField == "" {
		return "", fmt.Errorf("recieved empty API url from response: \n %s", string(bodyBytes))
	}

	if forkResponse.HtmlUrlField == "" {
		return "", fmt.Errorf("recieved empty Git url from response: \n %s", string(bodyBytes))
	}

	return forkResponse.HtmlUrlField, nil
}
