package services

import (
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v33/github"
)

var (
	GitServiceObject gitServiceInterface = gitServiceImplementation{}
)

type gitServiceInterface interface {
	getGitHubClient() *github.Client
	CheckUserAccessRepo(owner string, repo string) (*github.Repository, error)
	CloneRepo(owner string, repo string) (*git.Repository, string, error)
	CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error)
	CreateSecretFile(path string, secretFile string) error
	EditSecretFile(path string, secretsChanges secretUpdateMap) error
	CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error
	ForkRepo(owner string, repo string) (forkedOwner interface{}, gitURL interface{}, err error)
}

type gitServiceImplementation struct { }
