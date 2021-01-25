package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// Background: Third Party packages.
var (
	ContextBackground       = context.Background
	Oauth2StaticTokenSource = oauth2.StaticTokenSource
	Oauth2NewClient         = oauth2.NewClient
	GithubNewClient         = github.NewClient
	GithubRepositories      = GitServiceObject.GetGitHubClient().Repositories.Get
	HTTPGetCheckForkedRepo  = http.Get
	HTTPPostForkRepo        = http.PostForm
	IoutilReadAll           = ioutil.ReadAll
	JSONUnmarshal           = json.Unmarshal
	OsStat                  = os.Stat
	OsIsNotExist            = os.IsNotExist
	IoutilWriteFile         = ioutil.WriteFile
	GitPlainClone           = git.PlainClone
	OsRemoveAll             = os.RemoveAll
	IoutilReadFile          = ioutil.ReadFile
	ReflectValueOf          = reflect.ValueOf
	JSONMarshalIndent       = json.MarshalIndent
)

type gitServiceInterface interface {
	GetGitHubClient() *github.Client
	CheckUserAccessRepo(owner string, repo string) (*github.Repository, error)
	CloneRepo(string, string) (*git.Repository, string, error)
	CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error)
	CreateSecretFile(path string, secretFile string) error
	EditSecretFile(path string, secretsChanges SecretUpdateMap) error
	CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error
	ForkRepo(owner string, repo string) (forkedOwner interface{}, gitURL interface{}, err error)
	CheckForkedRepo(url string) error
}

type thirdPartyGitHubInterface interface {
	Head(*git.Repository) (*plumbing.Reference, error)
	Worktree(*git.Repository) (*git.Worktree, error)
	Fetch(*git.Repository) error
	Checkout(*git.Worktree, string, *plumbing.Reference) (string, error)
	Add(*git.Worktree, string) (plumbing.Hash, error)
	Commit(*git.Worktree, string, string) (plumbing.Hash, error)
	CommitObject(*git.Repository, plumbing.Hash) (*object.Commit, error)
	Push(*git.Repository) error
	PullRequest(context.Context, *github.Client, string, string, string, string, string, string, string) (*github.PullRequest, *github.Response, error)
}
