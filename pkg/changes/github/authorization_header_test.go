package github_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/changes/auth/authfakes"

	"github.com/calebamiles/keps/pkg/changes/github"
)

var _ = Describe("GitHub authorization", func() {
	Describe("adding a user token to an HTTP request", func() {
		It("adds a GitHub authorization header", func() {
			githubToken := "this-isnt-a-valid-token"

			token := &authfakes.FakeTokenProvider{}
                        token.ValueReturns(githubToken, nil)

			req, err := http.NewRequest(http.MethodGet, "https://kubernetes.io", nil)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a HTTP GET request with a valid URL")

			err = github.AddAuthorizationHeader(req, token)
			Expect(err).ToNot(HaveOccurred(), "expected no error when adding a GitHub authorization header to a valid HTTP request")

			tokenVal := req.Header.Get(github.AuthorizationHeaderName)
			Expect(tokenVal).To(ContainSubstring(githubToken), "expected authorization header to include given token")
		})
	})
})
