package services

import (
	"context"
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"net/http"
)

type contextMock struct {
	BackgoundHandler     			 func() context.Context
}

type oAuthMock struct {
	StaticTokenSourceHandler		 func(*oauth2.Token) oauth2.TokenSource
	OauthNewClientHandler			 func(context.Context, oauth2.TokenSource) *http.Client
}

type gitServiceMock struct {
	GithubNewClientHandler			 func(*http.Client) *github.Client
}

func (mock contextMock) Background() context.Context {
	return mock.BackgoundHandler()
}

func (mock oAuthMock) StaticTokenSource(token *oauth2.Token) oauth2.TokenSource {
	return mock.StaticTokenSourceHandler(token)
}

// client for OAuth
func (mock oAuthMock) NewClient(ctx context.Context, src oauth2.TokenSource) *http.Client {
	return mock.OauthNewClientHandler(ctx, src)
}

// client for Github
func (mock gitServiceMock) NewClient(client *http.Client) *github.Client {
	return mock.GithubNewClientHandler(client)
}

var _ = Describe("Services", func() {

	It("test", func() {
		contextObj := contextMock{}
		oAuthObj := oAuthMock{}
		gitServiceObj := gitServiceMock{}
		contextObj.BackgoundHandler = func() context.Context {
			var ctx context.Context
			return ctx
		}
		oAuthObj.StaticTokenSourceHandler = func(*oauth2.Token) oauth2.TokenSource {
			var src oauth2.TokenSource
			return src
		}
		oAuthObj.OauthNewClientHandler = func(context.Context, oauth2.TokenSource) *http.Client {
			return new(http.Client)
		}
		gitServiceObj.GithubNewClientHandler = func(*http.Client) *github.Client {
			return new(github.Client)
		}
		ThirdPartyContext = contextObj
		ThirdPartyOauth = oAuthObj
		ThirdPartyGitHub = gitServiceObj
		test := GitServiceObject.GetGitHubClient()
		Expect(test).To(Equal(new(github.Client)))
	})
})
