package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type contextMock struct {
	BackgoundHandler func() context.Context
}

type oAuthMock struct {
	StaticTokenSourceHandler func(*oauth2.Token) oauth2.TokenSource
	OauthNewClientHandler func(context.Context, oauth2.TokenSource) *http.Client
}

type gitServiceMock struct {
	GithubNewClientHandler func(*http.Client) *github.Client
	GetRepoInfoHandler func(*github.Client, context.Context, string, string) (*github.Repository, *github.Response, error)
	HeadHandler func(*git.Repository) (*plumbing.Reference, error)
	WorktreeHandler func(*git.Repository) (*git.Worktree, error)
	FetchHandler func(repoGit *git.Repository) error
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

// get git repo info
func (mock gitServiceMock) Get(client *github.Client, ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
	return mock.GetRepoInfoHandler(client, ctx, owner, repo)
}

// get branch head
func (mock gitServiceMock) Head(repo *git.Repository) (*plumbing.Reference, error) {
	return mock.HeadHandler(repo)
}

// get branch worktree
func (mock gitServiceMock) Worktree(repo *git.Repository) (*git.Worktree, error) {
	return mock.WorktreeHandler(repo)
}

func (mock gitServiceMock) Fetch(repo *git.Repository) error {
	return mock.FetchHandler(repo)
}

var _ = Describe("Services", func() {
	Context(" when trying to create new github client", func() {
		It("verifies user by token and returns new github client", func() {
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
			client := GitServiceObject.GetGitHubClient()
			Expect(client).To(Equal(new(github.Client)))
		})
	})

	Context("when given user access token", func(){
		It("returns repo info", func(){
			gitServiceObj := gitServiceMock{}

			gitServiceObj.GithubNewClientHandler = func(*http.Client) *github.Client {
				return new(github.Client)
			}
			gitServiceObj.GetRepoInfoHandler = func(*github.Client, context.Context, string, string) (*github.Repository, *github.Response, error){
				return new(github.Repository), nil, nil
			}

			ThirdPartyGitHub = gitServiceObj
			result, err:= GitServiceObject.CheckUserAccessRepo("john", "repo")
			Expect(err).To(BeNil())
			Expect(result).To(Equal(new(github.Repository)))
		})

		It("returns error and exits method when problem occurs", func(){
			gitServiceObj := gitServiceMock{}

			gitServiceObj.GithubNewClientHandler = func(*http.Client) *github.Client {
				return new(github.Client)
			}
			gitServiceObj.GetRepoInfoHandler = func(*github.Client, context.Context, string, string) (*github.Repository, *github.Response, error){
				return nil, nil, errors.New("error while fetching repo")
			}

			ThirdPartyGitHub = gitServiceObj
			result, err:= GitServiceObject.CheckUserAccessRepo("john", "repo")
			Expect(strings.Contains(fmt.Sprintf("%v",err), "error")).To(BeTrue())
			Expect(result).To(BeNil())
		})
	})

	Context("when providing repo details", func(){
		It("returns current branch and head branch", func(){
			gitServiceObj := gitServiceMock{}

			gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
				return new(plumbing.Reference), nil
			}
			gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
				return new(git.Worktree), nil
			}
			gitServiceObj.FetchHandler = func(repo *git.Repository) error {
				return nil
			}
			ThirdPartyGitHub = gitServiceObj
			branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
			Expect(err).To(BeNil())
			Expect(branch).To(Equal(""))
			Expect(head).To(Equal(""))
		})

		It("returns error when problem occurs in fetching head", func(){
			gitServiceObj := gitServiceMock{}

			gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
				return nil , errors.New("error creating branch")
			}

			ThirdPartyGitHub = gitServiceObj
			branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
			Expect(strings.Contains(fmt.Sprintf("%v",err), "error")).To(BeTrue())
			Expect(branch).To(Equal(""))
			Expect(head).To(Equal(""))
		})

		It("returns error when problem occurs in fetching working branch", func(){
			gitServiceObj := gitServiceMock{}

			gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
				return new(plumbing.Reference), nil
			}

			gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
				return nil, errors.New("error fetching working branch")
			}

			ThirdPartyGitHub = gitServiceObj
			branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
			Expect(strings.Contains(fmt.Sprintf("%v",err), "error fetching working branch")).To(BeTrue())
			Expect(branch).To(Equal(""))
			Expect(head).To(Equal(""))
		})

		It("returns error when problem occurs in fetching all branches", func(){
			gitServiceObj := gitServiceMock{}

			gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
				return new(plumbing.Reference), nil
			}

			gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
				return nil, errors.New("error fetching all branches")
			}

			ThirdPartyGitHub = gitServiceObj
			branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
			Expect(strings.Contains(fmt.Sprintf("%v",err), "error fetching working branch")).To(BeTrue())
			Expect(branch).To(Equal(""))
			Expect(head).To(Equal(""))
		})
	})

})
