package main

import (
	"fmt"
)

type controllerInterface interface {
	UpdateSecretFile(data *UpdateParams) (int, string)
	CreateSecretFile(data *CreateParams) (int, string)
}

// ControllerObject: Exposing controllers to external packages (Ie, Unit tests)
var (
	ControllerObject controllerInterface = controllerImplementation{}
)

type controllerImplementation struct{}

func (controller controllerImplementation) CreateSecretFile(data *CreateParams) (int, string) {
	originalRepoURL := data.Repo
	originalOwner := data.Owner
	zeroLogger.Info().Msgf("REPO: %s", originalRepoURL)
	zeroLogger.Info().Msgf("OWNER: %s", originalOwner)

	_, err := GitServiceObject.CheckUserAccessRepo(originalOwner, originalRepoURL)
	if err != nil {
		zeroLogger.Error().Msgf("access denied: %v", err)
		return 403, "You do not have access to the repo"
	}

	_forkOwner, _, err := GitServiceObject.ForkRepo(originalOwner, originalRepoURL)
	if err != nil {
		zeroLogger.Error().Msgf("Fork Error: %v", err)
		return 400, fmt.Sprintf("Error Forking Repo: %v", err)
	}

	getURL := fmt.Sprintf("https://%s@api.github.com/repos/%s/%s", gitHubToken, _forkOwner, originalRepoURL)

	err = GitServiceObject.CheckForkedRepo(getURL)

	if err != nil {
		errorMsg := fmt.Sprintf("Repo didn't fork properly: %v", err)
		zeroLogger.Error().Msg(errorMsg)
		return 400, errorMsg
	}

	forkOwner := fmt.Sprintf("%v", _forkOwner)
	zeroLogger.Info().Msgf("Owner who forked the repo: %s", forkOwner)

	forkedRepoURL, path, err := GitServiceObject.CloneRepo(forkOwner, originalRepoURL)
	if err != nil {
		zeroLogger.Error().Msgf("Error: %v", err)
		return 400, fmt.Sprintf("Error Cloning Repo: %v", err)
	}

	currentBranch, headBranch, err := GitServiceObject.CreateBranchRepo(forkedRepoURL, originalRepoURL, "create")
	if err != nil {
		zeroLogger.Error().Msgf("Branch not created: %v", err)
		return 400, fmt.Sprintf("Error Creating Branch: %s", err)
	}

	err = GitServiceObject.CreateSecretFile(path, data.Content)
	if err != nil {
		zeroLogger.Error().Msgf("Secrets file not created: %v", err)
		return 400, fmt.Sprintf("Error creating %s file: %v", SecretsFileName, err)
	}

	var description = "Created and added .secrets.baseline file, the bot ran the scan on the whole repo " +
		"and found all the secrets and placed them in .secrets.baseline file."
	err = GitServiceObject.CreateCommitAndPr(forkOwner, originalOwner, originalRepoURL, currentBranch, headBranch, "Create", description, forkedRepoURL)
	if err != nil {
		zeroLogger.Info().Msg("Updated the existing PR")
		return 200, fmt.Sprintf("PR was Updated !")
	}
	zeroLogger.Info().Msg("PR was Created Successfully!")
	return 200, "PR was Created !"

}

func (controller controllerImplementation) UpdateSecretFile(data *UpdateParams) (int, string) {
	originalRepoURL := data.Repo
	originalOwner := data.Owner
	_, err := GitServiceObject.CheckUserAccessRepo(originalOwner, originalRepoURL)
	if err != nil {
		zeroLogger.Error().Msgf("access denied: %v", err)
		return 403, "You do not have access to the repo"
	}

	_forkOwner, _, err := GitServiceObject.ForkRepo(originalOwner, originalRepoURL)
	if err != nil {
		zeroLogger.Error().Msgf("Fork Error: %v", err)
		return 400, fmt.Sprintf("Error Forking Repo: %v", err)
	}

	getURL := fmt.Sprintf("https://%s@api.github.com/repos/%s/%s", gitHubToken, _forkOwner, originalRepoURL)

	err = GitServiceObject.CheckForkedRepo(getURL)

	if err != nil {
		errorMsg := fmt.Sprintf("Repo didn't fork properly: %v", err)
		zeroLogger.Error().Msg(errorMsg)
		return 400, errorMsg
	}

	forkOwner := fmt.Sprintf("%v", _forkOwner)
	zeroLogger.Info().Msgf("Owner who forked the repo: %s", forkOwner)

	forkedRepoURL, path, err := GitServiceObject.CloneRepo(forkOwner, originalRepoURL)
	if err != nil {
		zeroLogger.Error().Msgf("Error: %v", err)
		return 400, fmt.Sprintf("Error Cloning Repo: %v", err)
	}

	currentBranch, headBranch, err := GitServiceObject.CreateBranchRepo(forkedRepoURL, originalRepoURL, "create")
	if err != nil {
		zeroLogger.Error().Msgf("Branch not created: %v", err)
		return 400, fmt.Sprintf("Error Creating Branch: %s", err)
	}

	err = GitServiceObject.EditSecretFile(path, data.Changes)
	if err != nil {
		zeroLogger.Error().Msgf("Error editing the %s file: %v", SecretsFileName, err)
		return 400, fmt.Sprintf("Cannot edit %s file", SecretsFileName)
	}

	var description = "Updated .secrets.baseline file, the user marked the secrets as false positive and " +
		"sent those changes to the repo."
	err = GitServiceObject.CreateCommitAndPr(forkOwner, originalOwner, originalRepoURL, currentBranch, headBranch, "Update", description, forkedRepoURL)
	if err != nil {
		zeroLogger.Info().Msg("Updated the existing PR")
		return 200, fmt.Sprintf("PR was Updated !")
	}
	zeroLogger.Info().Msg("PR was Created Successfully!")
	return 200, "PR was Created !"

}
