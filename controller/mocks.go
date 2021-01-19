package controller

import (
	. "github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v33/github"
)

type gitServiceMock struct {
	GetGitHubClientHandler     func() *github.Client
	CheckUserAccessRepoHandler func(string, string) (*github.Repository, error)
	ForkRepoHandler            func(string, string) (interface{}, interface{}, error)
	CloneRepoHandler           func(string, string) (*git.Repository, string, error)
	CreateBranchRepoHandler    func(*git.Repository, string, string) (string, string, error)
	CreateSecretFileHandler    func(string, string) error
	CreateCommitAndPrHandler   func(string, string, string, string, string, string, string, *git.Repository) error
	EditSecretFileHandler      func(string, SecretUpdateMap) error
	CheckForkedRepoHandler	   func(string) error
}

type contextMock struct {
	BodyParserHandler func(*createParams) error
	StatusHandler     func(int) *fiber.Ctx
}

func (mock contextMock) BodyParser(data *createParams) error {
	return mock.BodyParserHandler(data)
}

func (mock contextMock) Status(code int) *fiber.Ctx {
	return mock.StatusHandler(code)
}

func (mock gitServiceMock) GetGitHubClient() *github.Client {
	return mock.GetGitHubClientHandler()
}

func (mock gitServiceMock) CheckUserAccessRepo(owner string, repo string) (*github.Repository, error) {
	return mock.CheckUserAccessRepoHandler(owner, repo)
}

func (mock gitServiceMock) ForkRepo(owner string, repo string) (interface{}, interface{}, error) {
	return mock.ForkRepoHandler(owner, repo)
}

func (mock gitServiceMock) CloneRepo(owner string, repo string) (*git.Repository, string, error) {
	return mock.CloneRepoHandler(owner, repo)
}

func (mock gitServiceMock) CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error) {
	return mock.CreateBranchRepoHandler(repoGit, repoName, action)
}

func (mock gitServiceMock) CreateSecretFile(path string, secretFile string) error {
	return mock.CreateSecretFileHandler(path, secretFile)
}

func (mock gitServiceMock) CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error {
	return mock.CreateCommitAndPrHandler(owner, originalOwner, repo, currentBranch, headBranch, action, description, repoGit)
}

func (mock gitServiceMock) EditSecretFile(path string, secretsChanges SecretUpdateMap) error {
	return mock.EditSecretFileHandler(path, secretsChanges)
}

func (mock gitServiceMock) CheckForkedRepo(getURL string) error {
	return mock.CheckForkedRepoHandler(getURL)
}