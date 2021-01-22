package services

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"
)

var (
	//GitServiceObject : Object with service methods.
	GitServiceObject gitServiceInterface = gitServiceImplementation{}
	ThirdPartyGitHub thirdPartyGitHubInterface = thirdPartyGitHubImpl{}
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
	ioutilReadFile          = ioutil.ReadFile
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
	Checkout(*git.Worktree, string, *plumbing.Reference) (error,bool)
	Add(*git.Worktree, string) (plumbing.Hash, error)
	Commit(*git.Worktree, string, string) (plumbing.Hash, error)
	CommitObject(*git.Repository, plumbing.Hash) (*object.Commit, error)
	Push(*git.Repository) error
	PullRequest(context.Context, *github.Client, string, string, string, string, string, string, string) (*github.PullRequest, *github.Response, error)
}


type gitServiceImplementation struct {}
type thirdPartyGitHubImpl struct {}

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

func (service thirdPartyGitHubImpl) Checkout(workingBranch *git.Worktree, branch string, headRef *plumbing.Reference) (error, bool) {
	var branchExists bool
	err :=  workingBranch.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err == nil {
		branchExists =  true
	}
	if err != nil {
		ZeroLogger.Info().Msgf("Creating new branch")

		err = workingBranch.Checkout(&git.CheckoutOptions{
			Hash:   headRef.Hash(),
			Branch: plumbing.NewBranchReferenceName(branch),
			Create: true,
		})
		if err != nil {
			branchExists = false
		}
		ZeroLogger.Info().Msg("Branch created!")
	}
	return err, branchExists
}

func (service thirdPartyGitHubImpl) Add(workingBranch *git.Worktree, SecretFile string) (plumbing.Hash, error) {
	return workingBranch.Add(SecretFile)
}

func (service thirdPartyGitHubImpl) Commit(workingBranch *git.Worktree, msg string, owner string) (plumbing.Hash, error) {
	return workingBranch.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name: owner,
			When: time.Now(),
		},
	})
}

func (service thirdPartyGitHubImpl) CommitObject(repoGit *git.Repository, commit plumbing.Hash) (*object.Commit, error) {
	return repoGit.CommitObject(commit)
}

func (service thirdPartyGitHubImpl) Push(repoGit *git.Repository) error {
	obj := &git.PushOptions{}
	return repoGit.Push(obj)
}

func (service thirdPartyGitHubImpl) PullRequest(
	ctx context.Context,
	client *github.Client,
	originalOwner string,
	repo string,
	action string,
	owner string,
	currentBranch string,
	headBranch string,
	description string) (*github.PullRequest, *github.Response, error) {
	newPR := &github.NewPullRequest{
		Title:               github.String(fmt.Sprintf("[Detect Secrets] %s Secret BaseLine File", action)),
		Head:                github.String(fmt.Sprintf("%s:%s", owner, currentBranch)),
		Base:                github.String(headBranch),
		Body:                github.String(description),
		MaintainerCanModify: github.Bool(true),
	}

	return client.PullRequests.Create(ctx, originalOwner, repo, newPR)
}
