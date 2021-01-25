package main_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eliezer-borde-globant/EBGoProject/main"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {
	Describe("Create Controller", func() {
		Context("CreateSecretFile controller is triggered", func() {
			gitService := gitServiceMock{}
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
			It("should create secrets file, and commit, push, and create PR", func() {
				main.GitServiceObject = gitService
				statusCode, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
				Expect(statusCode).To(Equal(200))
				Expect(msg).To(Equal("PR was Created !"))
			})

			Context("user doesn't have access to the repo", func() {
				It("should return the user access error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return nil, errors.New("error in checkUserAccess service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(statusCode).To(Equal(403))
					Expect(strings.Contains(msg, "You do not have access to the repo")).To(BeTrue())
				})
			})

			Context("problem occurs while forking the repo", func() {
				It("should return the fork error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "", "", errors.New("error in forkRepo service")
					}
					main.GitServiceObject = gitService
					status, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(status).To(Equal(400))
					Expect(strings.Contains(msg, "Error Forking Repo")).To(BeTrue())
				})
			})

			Context("check for forked repo fails", func() {
				It("should return the error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return errors.New("error in checkRepo service")
					}
					main.GitServiceObject = gitService
					status, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(status).To(Equal(400))
					Expect(strings.Contains(msg, "Repo didn't fork properly")).To(BeTrue())
				})
			})

			Context("there is problem in cloning the repo", func() {
				It("should return the clone error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return nil, "", errors.New("error in cloneRepo service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(statusCode).To(Equal(400))
					Expect(strings.Contains(msg, "Error Cloning Repo")).To(BeTrue())
				})
			})
			Context("problem occurs in creating branch", func() {
				It("should return error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return new(git.Repository), "path", nil
					}
					gitService.CreateBranchRepoHandler = func(*git.Repository, string, string) (string, string, error) {
						return "", "", errors.New("error in createBranch service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(statusCode).To(Equal(400))
					Expect(strings.Contains(msg, "Error Creating Branch")).To(BeTrue())
				})
			})
			Context("problem occurs in creating secrets file", func() {
				It("should return error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return new(git.Repository), "path", nil
					}
					gitService.CreateBranchRepoHandler = func(*git.Repository, string, string) (string, string, error) {
						return "branch", "headBranch", nil
					}
					gitService.CreateSecretFileHandler = func(string, string) error {
						return errors.New("error in createSecretFile service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(statusCode).To(Equal(400))
					Expect(strings.Contains(msg, fmt.Sprintf("Error creating %s file", main.SecretsFileName))).To(BeTrue())
				})
			})
			Context("creating of PR fails", func() {
				It("should return error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
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
						return errors.New("error in creating PR service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.CreateSecretFile(new(main.CreateParams))
					Expect(statusCode).To(Equal(200))
					Expect(strings.Contains(msg, "PR was Updated !")).To(BeTrue())
				})
			})
		})
	})

	Describe("Update Controller", func() {
		Context("UpdatedSecretFile controller is triggered", func() {
			gitService := gitServiceMock{}

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
			gitService.EditSecretFileHandler = func(string, main.SecretUpdateMap) error {
				return nil
			}
			gitService.CreateCommitAndPrHandler = func(string, string, string, string, string, string, string, *git.Repository) error {
				return nil
			}
			gitService.CheckForkedRepoHandler = func(string) error {
				return nil
			}
			It("should update secrets file, and commit, push, and create PR", func() {
				main.GitServiceObject = gitService
				statusCode, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
				Expect(statusCode).To(Equal(200))
				Expect(msg).To(Equal("PR was Created !"))
			})

			Context("user doesn't have access to the repo", func() {
				It("should return the user access error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return nil, errors.New("error in checkUserAccess service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(statusCode).To(Equal(403))
					Expect(strings.Contains(msg, "You do not have access to the repo")).To(BeTrue())
				})
			})

			Context("problem occurs while forking the repo", func() {
				It("should return the fork error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "", "", errors.New("error in forkRepo service")
					}
					main.GitServiceObject = gitService
					status, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(status).To(Equal(400))
					Expect(strings.Contains(msg, "Error Forking Repo")).To(BeTrue())
				})
			})

			Context("check for forked repo fails", func() {
				It("should return the error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return errors.New("error in checkRepo service")
					}
					main.GitServiceObject = gitService
					status, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(status).To(Equal(400))
					Expect(strings.Contains(msg, "Repo didn't fork properly")).To(BeTrue())
				})
			})

			Context("there is problem in cloning the repo", func() {
				It("should return the clone error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return nil, "", errors.New("error in cloneRepo service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(statusCode).To(Equal(400))
					Expect(strings.Contains(msg, "Error Cloning Repo")).To(BeTrue())
				})
			})
			Context("problem occurs in creating branch", func() {
				It("should return error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return new(git.Repository), "path", nil
					}
					gitService.CreateBranchRepoHandler = func(*git.Repository, string, string) (string, string, error) {
						return "", "", errors.New("error in createBranch service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(statusCode).To(Equal(400))
					Expect(strings.Contains(msg, "Error Creating Branch")).To(BeTrue())
				})
			})
			Context("problem occurs in editing secrets file", func() {
				It("should return error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return new(git.Repository), "path", nil
					}
					gitService.CreateBranchRepoHandler = func(*git.Repository, string, string) (string, string, error) {
						return "branch", "headBranch", nil
					}
					gitService.EditSecretFileHandler = func(string, main.SecretUpdateMap) error {
						return errors.New("error in editSecretFile service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(statusCode).To(Equal(400))
					Expect(strings.Contains(msg, fmt.Sprintf("Cannot edit %s file", main.SecretsFileName))).To(BeTrue())
				})
			})
			Context("creating of PR fails", func() {
				It("should return error", func() {
					gitService.CheckUserAccessRepoHandler = func(string, string) (*github.Repository, error) {
						return new(github.Repository), nil
					}
					gitService.ForkRepoHandler = func(string, string) (interface{}, interface{}, error) {
						return "username", "http://github.com/username/test", nil
					}
					gitService.CheckForkedRepoHandler = func(string) error {
						return nil
					}
					gitService.CloneRepoHandler = func(string, string) (*git.Repository, string, error) {
						return new(git.Repository), "path", nil
					}
					gitService.CreateBranchRepoHandler = func(*git.Repository, string, string) (string, string, error) {
						return "branch", "headBranch", nil
					}
					gitService.EditSecretFileHandler = func(string, main.SecretUpdateMap) error {
						return nil
					}
					gitService.CreateCommitAndPrHandler = func(string, string, string, string, string, string, string, *git.Repository) error {
						return errors.New("error in creating PR service")
					}
					main.GitServiceObject = gitService
					statusCode, msg := main.ControllerObject.UpdateSecretFile(new(main.UpdateParams))
					Expect(statusCode).To(Equal(200))
					Expect(strings.Contains(msg, "PR was Updated !")).To(BeTrue())
				})
			})
		})
	})

})
