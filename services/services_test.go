package services

import (
	"context"
	"net/http"

	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("Services", func() {

	Context("problem occurs while forking the repo", func() {
		It("should return the fork error", func() {
			BackgroundMock := ContextBackground
			defer func() { ContextBackground = BackgroundMock }()
			ContextBackground = func() context.Context {
				var ctx context.Context
				return ctx
			}
			StaticTokenSourceMock := Oauth2StaticTokenSource
			defer func() { Oauth2StaticTokenSource = StaticTokenSourceMock }()
			Oauth2StaticTokenSource = func(*oauth2.Token) oauth2.TokenSource {
				var src oauth2.TokenSource
				return src
			}
			NewClientMock := Oauth2NewClient
			defer func() { Oauth2NewClient = NewClientMock }()
			Oauth2NewClient = func(context.Context, oauth2.TokenSource) *http.Client {
				return new(http.Client)
			}
			GitNewClientMock := GithubNewClient
			defer func() { GithubNewClient = GitNewClientMock }()
			GithubNewClient = func(*http.Client) *github.Client {
				return new(github.Client)
			}
			result := GitServiceObject.GetGitHubClient()
			Expect(result).To(Equal(new(github.Client)))

		})
		It("should return the fork error", func() {

			GithubRepositoriesMock := GithubRepositories
			defer func() { GithubRepositories = GithubRepositoriesMock }()
			GithubRepositories = func(context.Context, string, string) (*github.Repository, *github.Response, error) {
				return new(github.Repository), new(github.Response), nil
			}
			result, error := GitServiceObject.CheckUserAccessRepo("owner_test", "repo_test")
			Expect(error).To(BeNil())
			Expect(result).To(Equal(new(github.Repository)))

		})
	})

})
