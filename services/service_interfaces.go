package services

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/oauth2"
	//. "github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/google/go-github/v33/github"
	"net/http"
)

var (
	GitServiceObject gitServiceInterface = gitServiceImplementation{}
	ThirdPartyContext thirdPartyContextInterface = thirdPartyContextImpl{}
	ThirdPartyOauth thirdPartyOauthInterface = thirdPartyOauthImpl{}
	ThirdPartyGitHub thirdPartyGitHubInterface = thirdPartyGitHubImpl{}
)

type gitServiceInterface interface {
	GetGitHubClient() *github.Client
	CheckUserAccessRepo(owner string, repo string) (*github.Repository, error)
	//CloneRepo(owner string, repo string) (*git.Repository, string, error)
	CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error)
	//CreateSecretFile(path string, secretFile string) error
	//EditSecretFile(path string, secretsChanges SecretUpdateMap ) error
	//CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error
	//ForkRepo(owner string, repo string) (forkedOwner interface{}, gitURL interface{}, err error)
	//CheckForkedRepo(url string) error
}

type thirdPartyContextInterface interface {
	Background() context.Context
}

type thirdPartyOauthInterface interface {
	StaticTokenSource(*oauth2.Token) oauth2.TokenSource
	NewClient(context.Context, oauth2.TokenSource) *http.Client
}

type thirdPartyGitHubInterface interface {
	NewClient(*http.Client) *github.Client
	Get(*github.Client, context.Context, string, string) (*github.Repository, *github.Response, error)
	Head(*git.Repository) (*plumbing.Reference, error)
	Worktree(*git.Repository) (*git.Worktree, error)
	Fetch(*git.Repository) error
}


type gitServiceImplementation struct {}
type thirdPartyContextImpl struct {}
type thirdPartyOauthImpl struct {}
type thirdPartyGitHubImpl struct {}


func (service thirdPartyContextImpl) Background() context.Context {
	return context.Background()
}

func (service thirdPartyOauthImpl) StaticTokenSource(t *oauth2.Token) oauth2.TokenSource {
	return oauth2.StaticTokenSource(t)
}

func (service thirdPartyOauthImpl) NewClient(ctx context.Context, src oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, src)
}

func (service thirdPartyGitHubImpl) NewClient(httpClient *http.Client) *github.Client {
	return github.NewClient(httpClient)
}

func (service thirdPartyGitHubImpl) Get(client *github.Client, ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
	return client.Repositories.Get(ctx, owner, repo)
}

func (service thirdPartyGitHubImpl) Head(repoGit *git.Repository) (*plumbing.Reference, error) {
	return repoGit.Head()
}

func (service thirdPartyGitHubImpl) Worktree(repoGit *git.Repository) (*git.Worktree, error) {
	return repoGit.Worktree()
}

func (service thirdPartyGitHubImpl) Fetch(repoGit *git.Repository) error {
	return repoGit.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
}