package services

import (
	"github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v33/github"
)

var (
	//GitServiceObject : Object with service methods.
	GitServiceObject gitServiceInterface = gitServiceImplementation{}
)

type gitServiceInterface interface {
	GetGitHubClient() *github.Client
	CheckUserAccessRepo(owner string, repo string) (*github.Repository, error)
	CloneRepo(string, string) (*git.Repository, string, error)
	CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error)
	CreateSecretFile(path string, secretFile string) error
	EditSecretFile(path string, secretsChanges utils.SecretUpdateMap) error
	CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error
	ForkRepo(owner string, repo string) (forkedOwner interface{}, gitURL interface{}, err error)
	CheckForkedRepo(url string) error
}

type gitServiceImplementation struct{}
