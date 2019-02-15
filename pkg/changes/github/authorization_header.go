package github

import (
	"fmt"
	"net/http"

	"github.com/calebamiles/keps/pkg/changes/auth"
)

func AddAuthorizationHeader(req *http.Request, token auth.TokenProvider) error {
	authToken, err := token.Value()
	if err != nil {
		return err
	}

	req.Header.Add(AuthorizationHeaderName, fmt.Sprintf("token %s", authToken))

	return nil
}

const (
	AuthorizationHeaderName = "Authorization"
)
