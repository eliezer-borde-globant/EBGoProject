package controller

import (
	. "github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type gitServiceMock struct {
	CheckUserAccessRepoHandler func(string, string) (*github.Repository, error)
	ForkRepoHandler func(string, string) (interface{}, interface{}, error)
	CloneRepoHandler func(string, string) (*git.Repository, string, error)
	CreateBranchRepoHandler func(*git.Repository, string, string) (string, string, error)
	CreateSecretFileHandler func(string, string) error
	CreateCommitAndPrHandler func(string, string, string, string, string, string, string, *git.Repository) error
	EditSecretFileHandler func(string, SecretUpdateMap) error
}

type contextMock struct {
	BodyParserHandler func(*createParams) error
	StatusHandler func(int) *fiber.Ctx
}

func (mock contextMock) BodyParser(data *createParams) error {
	return mock.BodyParserHandler(data)
}

func (mock contextMock) Status(code int) *fiber.Ctx {
	return mock.StatusHandler(code)
}

func (mock gitServiceMock) CheckUserAccessRepo(owner string, repo string) (*github.Repository, error) {
	return mock.CheckUserAccessRepoHandler(owner , repo)
}

func (mock gitServiceMock) ForkRepo(owner string, repo string) (interface{}, interface{}, error) {
	return mock.ForkRepoHandler(owner , repo)
}

func (mock gitServiceMock) CloneRepo(owner string, repo string) (*git.Repository, string, error) {
	return mock.CloneRepoHandler(owner , repo)
}

func (mock gitServiceMock) CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error) {
	return mock.CreateBranchRepoHandler(repoGit , repoName, action)
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

var _ = Describe("Controller", func() {
	When("CreateSecretFile is triggered", func() {
		It("should create secrets file, and commit, push, and create PR", func() {
			gitService := gitServiceMock{}
			context := contextMock{}

			gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
				return new(github.Repository), nil
			}

			gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
				return "username", "http://github.com/username/test", nil
			}

			gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
				return new(git.Repository), "path", nil
			}

			gitService.CreateBranchRepoHandler = func(*git.Repository, string, string) (string, string, error) {
				return "branch", "headBranch", nil
			}

			gitService.CreateSecretFileHandler = func(string, string) error {
				return nil
			}

			gitService.CreateCommitAndPrHandler = func(string, string, string, string, string, string, string, *git.Repository) error {
				return nil
			}

			context.BodyParserHandler = func(*createParams) error {
				return nil
			}

			context.StatusHandler = func(code int) *fiber.Ctx {
				return new(fiber.Ctx)
			}

			result := ControllerObject.CreateSecretFile(context)

			Expect(result).To(Equal(nil))

		})
	})
})