package services

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	//. "github.com/eliezer-borde-globant/EBGoProject/utils"
	//"github.com/go-git/go-git/v5"
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
	//CheckUserAccessRepo(owner string, repo string) (*github.Repository, error)
	//CloneRepo(owner string, repo string) (*git.Repository, string, error)
	//CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error)
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
	fmt.Print("How are you")
	return oauth2.NewClient(ctx, src)
}

func (service thirdPartyGitHubImpl) NewClient(httpClient *http.Client) *github.Client {
	fmt.Println("hello there")
	return github.NewClient(httpClient)
}