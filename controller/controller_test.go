package controller

import (
	"errors"
	"github.com/eliezer-borde-globant/EBGoProject/services"
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)



var _ = Describe("Controller", func() {
	When("CreateSecretFile is triggered", func() {
		gitService := gitServiceMock{}
		context := contextMock{}
		gitService.GetGitHubClientHandler = func() *github.Client {
			return new(github.Client)
		}
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
		gitService.CheckForkedRepoHandler = func(string) error {
			return nil
		}
		context.BodyParserHandler = func(*createParams) error {
			return nil
		}
		context.StatusHandler = func(code int) *fiber.Ctx {
			return new(fiber.Ctx)
		}
		It("should create secrets file, and commit, push, and create PR", func() {
			services.GitServiceObject = gitService
			statusCode, msg := ControllerObject.CreateSecretFile(context)
			Expect(statusCode).To(Equal(200))
			Expect(msg).To(Equal("PR was Created !"))
		})
		It("test",  func() {
			gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
				return new(git.Repository), "path", errors.New("error cloning repo")
			}
			services.GitServiceObject = gitService
			statusCode, msg := ControllerObject.CreateSecretFile(context)
			Expect(statusCode).To(Equal(400))
			Expect(strings.Contains(msg, "Error Cloning Repo")).To(Equal(true))
		})
	})
})