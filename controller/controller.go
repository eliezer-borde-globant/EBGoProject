package controller

import (
	"fmt"
	. "github.com/eliezer-borde-globant/EBGoProject/services"
	. "github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
)


func (controller controllerImplementation) CreateSecretFile(c contextCreateInterface) error {
	data := new(createParams)
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error in data, please review input data: %s", err))
	}
	originalRepoURL := data.Repo
	originalOwner := data.Owner
	ZeroLogger.Info().Msgf("REPO: %s", originalRepoURL)
	ZeroLogger.Info().Msgf("OWNER: %s", originalOwner)

	_, err := GitServiceObject.CheckUserAccessRepo(originalOwner, originalRepoURL)
	if err != nil {
		ZeroLogger.Fatal().Msgf("%v", err)
		return c.Status(403).SendString("You do not have access to the repo")
	}

	_forkOwner, _, err := GitServiceObject.ForkRepo(originalOwner, originalRepoURL)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error Forking Repo: %s", err))
	}

	ZeroLogger.Info().Msgf("Checking if repo was forked properly")
	getURL := fmt.Sprintf("https://%s@api.github.com/repos/%s/%s", GitHubToken, _forkOwner, originalRepoURL)
	for {
		response, err := http.Get(getURL)
		if response.StatusCode == 200 {
			ZeroLogger.Info().Msgf("Repo has been forked successfully for user: %s", _forkOwner)
			break
		}
		if err != nil {
			panic(err)
		}
	}

	forkOwner := fmt.Sprintf("%v", _forkOwner)
	ZeroLogger.Info().Msgf("Owner who forked the repo: %s", forkOwner)

	forkedRepoURL, path, err := GitServiceObject.CloneRepo(forkOwner, originalRepoURL)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error Cloning Repo: %s", err))
	}

	currentBranch, headBranch, err := GitServiceObject.CreateBranchRepo(forkedRepoURL, originalRepoURL, "create")
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error Creating Branch: %s", err))
	}

	err = GitServiceObject.CreateSecretFile(path, data.Content)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error creating .secrets.baseline file: %s", err))
	}

	var description = "Created and added .secrets.baseline file, the bot ran the scan on the whole repo " +
		"and found all the secrets and placed them in .secrets.baseline file."
	err = GitServiceObject.CreateCommitAndPr(forkOwner, originalOwner, originalRepoURL, currentBranch, headBranch, "Create", description, forkedRepoURL)
	if err != nil {
		return c.Status(200).SendString(fmt.Sprintf("PR was Updated !"))
	}
	return c.Status(200).SendString(fmt.Sprintf("PR was Created !"))

}

func (controller controllerImplementation) UpdateSecretFile(c *fiber.Ctx) error {
	data := new(updateParams)
	if err := c.BodyParser(data); err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error in data, please review input data: %s", err))
	}
	originalRepoURL := data.Repo
	originalOwner := data.Owner
	_, err := GitServiceObject.CheckUserAccessRepo(originalOwner, originalRepoURL)
	if err != nil {
		ZeroLogger.Fatal().Msgf("%v", err)
		return c.Status(403).SendString("You do not have access to the repo")
	}
	_forkOwner, _, err := GitServiceObject.ForkRepo(originalOwner, originalRepoURL)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error Forking Repo: %s", err))
	}

	ZeroLogger.Info().Msgf("Checking if repo was forked properly")
	getURL := fmt.Sprintf("https://%s@api.github.com/repos/%s/%s", GitHubToken, _forkOwner, originalRepoURL)
	for {
		response, err := http.Get(getURL)
		if response.StatusCode == 200 {
			ZeroLogger.Info().Msgf("Repo has been forked successfully for user: %s", _forkOwner)
			break
		}
		if err != nil {
			panic(err)
		}
	}

	forkOwner := fmt.Sprintf("%v", _forkOwner)
	forkedRepoURL, path, err := GitServiceObject.CloneRepo(forkOwner, originalRepoURL)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error Cloning Repo: %s", err))
	}

	currentBranch, headBranch, err := GitServiceObject.CreateBranchRepo(forkedRepoURL, originalRepoURL, "update")
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error Creating Branch: %s", err))
	}
	err = GitServiceObject.EditSecretFile(path, data.Changes)
	if err != nil {
		ZeroLogger.Fatal().Msgf("%v", err)
		return c.Status(400).SendString("Cannot edit the secretfile")
	}
	var description = "Updated .secrets.baseline file, the user marked the secrets as false positive and " +
		"sent those changes to the repo."
	err = GitServiceObject.CreateCommitAndPr(forkOwner, originalOwner, originalRepoURL, currentBranch, headBranch, "Update", description, forkedRepoURL)
	if err != nil {
		return c.Status(200).SendString(fmt.Sprintf("PR was Updated !"))
	}
	return c.Status(200).SendString(fmt.Sprintf("PR was Created !"))

}