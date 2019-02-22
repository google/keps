package planctae

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	authorizationApiUrl = "https://api.github.com/user/memberships/orgs/planctae"

	activeMembership  = "active"
	pendingMembership = "pending"
)

// IsAuthorized determines whether the principal is allowed to
// interact with the `Planctae` GitHub organization
// https://developer.github.com/v3/orgs/members/#get-your-organization-membership
// TODO define an internal interface to use rather than requesting the full settings.Runtime
func IsAuthorized(runtime settings.Runtime) (bool, error) {
	token, err := github.NewTokenProvider(runtime.TokenPath())
	if err != nil {
		return false, err
	}

	authToken, err := token.Value()
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodGet, authorizationApiUrl, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add(github.AuthorizationHeaderName, github.AuthorizationHeaderValue(authToken))

	ctx, cancenl := context.WithTimeout(context.Background(), 5000*time.Milisecond)
	isAuthorized := req.WithContext(ctx)
	defer cancel()

	resp, err := github.HttpClient().Do(isAuthorized)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("expected status 200 OK, got %s when attempting to determine membership in %s", resp.Status, Organization)
	}

	var organizationMembershipResponse struct {
		State string `json:"state"`
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(bodyBytes, &organizationMembershipResponse)
	if err != nil {
		return false, err
	}

	switch oprganizationMembershipResponse.State {
	case "":
		return false, fmt.Errorf("empty reponse for membership state received from response:\n %s", string(bodyBytes))
	case pendingMembership:
		return false, errors.New("membership within Planctae is pending, please accept before continuing")
	case activeMembership:
		return true, nil
	}
}
