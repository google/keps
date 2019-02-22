package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/calebamiles/keps/pkg/changes/auth"
)

type PullRequestCreator func(token auth.TokenProvider, routingInfo PullRequestRoutingInfo, prTitle PullRequestTitle, prDescription PullRequestDescription) (string, error)

func CreatePullRequest(token auth.TokenProvider, routingInfo PullRequestRoutingInfo, prTitle PullRequestTitle, prDescription PullRequestDescription) (string, error) {
	// serialize request payload
	var createPrPayload struct {
		TitleField          PullRequestTitle       `json:"title"`
		DescriptionField    PullRequestDescription `json:"body,omitempty"`
		SourceBranchField   string                 `json:"head"`
		TargetBranchField   string                 `json:"base"`
		MaintainerCanModify bool                   `json:"maintainer_can_modify,omitempty"`
	}

	createPrPayload.TitleField = prTitle
	createPrPayload.DescriptionField = prDescription
	createPrPayload.TargetBranchField = routingInfo.TargetBranch()
	createPrPayload.SourceBranchField = fmt.Sprintf("%s:%s", routingInfo.SourceRepositoryOwner(), routingInfo.SourceBranch())
	createPrPayload.MaintainerCanModify = true

	payloadBytes, err := json.Marshal(createPrPayload)
	if err != nil {
		return "", err
	}

	payload := ioutil.NopCloser(bytes.NewBuffer(payloadBytes))

	remoteOwner := routingInfo.TargetRepositoryOwner()
	remoteRepo := routingInfo.TargetRepository()

	// create HTTP request
	req, err := http.NewRequest(http.MethodPost, PrUrl(remoteOwner, remoteRepo), payload)
	if err != nil {
		return "", err
	}

	// add auth header
	err = AddAuthorizationHeader(req, token)
	if err != nil {
		return "", err
	}

	// set context
	ctx, cancel := context.WithTimeout(context.Background(), 7000*time.Millisecond)
	createPullRequest := req.WithContext(ctx)
	defer cancel()

	c := HttpClient()

	// Do request
	resp, err := c.Do(createPullRequest)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("expected status code 201 Created, got: %s.\nURL: %s.\nBody: %s", resp.Status, PrUrl(remoteOwner, remoteRepo), string(payloadBytes))
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
