package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func (gitService gitServiceImplementation) GetGitHubClient() *github.Client {
	ctx := ContextBackground()
	ts := Oauth2StaticTokenSource(
		&oauth2.Token{AccessToken: gitHubToken},
	)
	tc := Oauth2NewClient(ctx, ts)
	return GithubNewClient(tc)
}

func (gitService gitServiceImplementation) CheckUserAccessRepo(owner string, repo string) (*github.Repository, error) {
	zeroLogger.Info().Msgf("check user has access to %s/%s", owner, repo)
	ctx := ContextBackground()
	repoInfo, _, err := GithubRepositories(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	fmt.Println(err)
	return repoInfo, nil
}

func (gitService gitServiceImplementation) CloneRepo(owner string, repo string) (*git.Repository, string, error) {
	path := fmt.Sprintf("/tmp/%s-%s", owner, repo)
	zeroLogger.Info().Msgf("Creating folder to clone %s", path)
	if _, err := os.Stat(path); !OsIsNotExist(err) {
		err := OsRemoveAll(path)
		if err != nil {
			zeroLogger.Error().Msgf("Error path to clone repo from %s/%s, error: %v", owner, repo, err)
			return nil, "", err
		}
	}
	zeroLogger.Info().Msg("Starting to clone Repo")
	repoInfo, err := GitPlainClone(path, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://%s@github.com/%s/%s", gitHubToken, owner, repo),
		Progress: os.Stdout,
	})
	if err != nil {
		zeroLogger.Error().Msgf("Error Cloning repo from %s/%s, error: %v", owner, repo, err)
		return nil, "", err
	}
	zeroLogger.Info().Msgf("Repo was cloned")
	return repoInfo, path, nil
}

func (gitService gitServiceImplementation) CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error) {
	zeroLogger.Info().Msgf("Creating Branch to update secret file in repo %s", repoName)
	headRef, err := ThirdPartyGitHub.Head(repoGit)
	if err != nil {
		zeroLogger.Error().Msgf("Error Creating Branch to update secret file in repo %s, error: %v", repoName, err)
		return "", "", err
	}
	headBranchName := strings.ReplaceAll(headRef.Name().String(), "refs/heads/", "")
	branch := fmt.Sprintf("secret_scanner_api/%s/%s/secrets_baseline_file", repoName, action)
	workingBranch, err := ThirdPartyGitHub.Worktree(repoGit)
	if err != nil {
		zeroLogger.Error().Msgf("Error Creating Branch to update secret file in repo %s, error: %v", repoName, err)
		return "", "", err
	}
	zeroLogger.Info().Msgf("Fetching all Branches from %s", repoName)
	err = ThirdPartyGitHub.Fetch(repoGit)
	if err != nil {
		zeroLogger.Error().Msgf("Error fetching remote Branches from repo %s, error: %v", repoName, err)
		return "", "", err
	}
	zeroLogger.Info().Msgf("Checking if the branch %s exists in %s", branch, repoName)
	branchStatus, err := ThirdPartyGitHub.Checkout(workingBranch, branch, headRef)
	if err != nil {
		zeroLogger.Error().Msgf("Something went wrong in checkout Method %s, error: %v", err)
		return "", "", err
	}
	if branchStatus == "EXISTING_BRANCH" {
		zeroLogger.Info().Msgf("Branch %s already exists in %s, Checking out...", branch, repoName)
	} else if branchStatus == "ERROR" {
		zeroLogger.Error().Msgf("Error Creating Branch to update secret file in repo %s, error: %v", repoName, err)
		return "", "", err
	} else if branchStatus == "NEW_BRANCH" {
		zeroLogger.Info().Msgf("Branch %s was created properly in ", branch, repoName)
	}
	return branch, headBranchName, err
}

func (gitService gitServiceImplementation) CreateSecretFile(path string, secretFile string) error {
	zeroLogger.Info().Msg(fmt.Sprintf("Creating Path %s to add %s file ", path, SecretsFileName))
	if _, err := OsStat(path); OsIsNotExist(err) {
		zeroLogger.Error().Msgf("*** Error Creating Path %s, error: %v", path, err)
		return err
	}
	path = fmt.Sprintf("%s/%s", path, SecretsFileName)
	err := IoutilWriteFile(path, []byte(secretFile), 0644)
	zeroLogger.Info().Msgf("File was created with the content at path: '%s'", path)
	if err != nil {
		zeroLogger.Error().Msgf("Error creating %s file in, %s error: %v", SecretsFileName, path, err)
		return err
	}
	return nil
}

func (gitService gitServiceImplementation) EditSecretFile(path string, secretsChanges SecretUpdateMap) error {
	zeroLogger.Info().Msgf("Starting to edit the secret file at path: '%s'", path)
	path = fmt.Sprintf("%s/%s", path, SecretsFileName)
	dat, err := IoutilReadFile(path)
	if err != nil {
		zeroLogger.Info().Msgf("Error..")
		return err
	}
	var fileStruct map[string]interface{}
	err = JSONUnmarshal(dat, &fileStruct)
	if err != nil {
		zeroLogger.Info().Msgf("Error getting json from existing file..")
		return err
	}
	results, ok := fileStruct["results"].(map[string]interface{})
	if !ok {
		err := errors.New("could not parse the result data in secret file, please check the data")
		zeroLogger.Error().Msgf("Error: %v", err)
		return err
	}

	for filename, secretData := range secretsChanges {
		_, ok := results[filename]
		if !ok {
			continue
		}
		for _, secret := range secretData {
			fileData := results[filename]
			fileSecrets := ReflectValueOf(fileData)
			for i := 0; i < fileSecrets.Len(); i++ {
				value := fileSecrets.Index(i)
				secrets := value.Interface().(map[string]interface{})
				if secret["hashed_secret"] == secrets["hashed_secret"] && secret["line_number"] == secrets["line_number"] {
					secrets["is_secret"] = secret["is_secret"]
					value.Set(ReflectValueOf(secrets))
				}
			}
			results[filename] = fileSecrets.Interface()
		}
	}
	fileStruct["results"] = results
	file, parseError := JSONMarshalIndent(fileStruct, "", "  ")
	if parseError != nil {
		zeroLogger.Error().Msgf("Cannot indent content of the file : %v", parseError)
		return parseError
	}
	writeFileError := IoutilWriteFile(path, file, 0644)
	if writeFileError != nil {
		zeroLogger.Error().Msgf("Error writing file: %v", parseError)
		return writeFileError
	}
	return nil
}

func (gitService gitServiceImplementation) CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error {
	zeroLogger.Info().Msg("Getting current branch")
	zeroLogger.Info().Msgf("Head Branch: %s", headBranch)
	zeroLogger.Info().Msgf("Current Branch: %s", currentBranch)
	workingBranch, err := ThirdPartyGitHub.Worktree(repoGit)
	if err != nil {
		zeroLogger.Info().Msgf("Error getting current branch '%s/%s'", owner, repo)
		return err
	}
	zeroLogger.Info().Msgf("Adding %s file to new branch ", SecretsFileName)
	_, err = ThirdPartyGitHub.Add(workingBranch, SecretsFileName)
	if err != nil {
		return err
	}
	zeroLogger.Info().Msgf("%s was added to stage ", SecretsFileName)
	zeroLogger.Info().Msg("Committing Changes")
	commit, err := ThirdPartyGitHub.Commit(workingBranch, fmt.Sprintf("chore: %s secret baseline file", action), owner)
	if err != nil {
		zeroLogger.Error().Msgf("Error Committing changes: %v", err)
		return err
	}
	zeroLogger.Info().Msg("Changes were committed")
	_, err = ThirdPartyGitHub.CommitObject(repoGit, commit)
	if err != nil {
		zeroLogger.Error().Msgf("Error Committing: %v", err)
		return err
	}
	zeroLogger.Info().Msgf("Commit created in '%s/%s'", owner, repo)

	zeroLogger.Info().Msg("Pushing changes to remote")
	err = ThirdPartyGitHub.Push(repoGit)
	if err != nil {
		zeroLogger.Info().Msg("Pushing changes to branch failed")
		return err
	}
	zeroLogger.Info().Msgf("Branch was pushed '%s/%s'", owner, repo)
	githubClient := GitServiceObject.GetGitHubClient()
	_, _, err = ThirdPartyGitHub.PullRequest(
		ContextBackground(), githubClient, originalOwner, repo, action,
		owner, currentBranch, headBranch, description)
	if err != nil {
		zeroLogger.Info().Msgf("PR success Updated! '%s/%s'", owner, repo)
		return err
	}
	zeroLogger.Info().Msgf("PR success Created! '%s/%s'", owner, repo)
	return nil
}

func (gitService gitServiceImplementation) ForkRepo(owner string, repo string) (forkedOwner interface{}, gitURL interface{}, err error) {
	zeroLogger.Info().Msgf("Forking repo from '%s/%s'", owner, repo)
	forkURL := fmt.Sprintf("https://%s@api.github.com/repos/%s/%s/forks", gitHubToken, owner, repo)
	resp, errPost := HTTPPostForkRepo(forkURL, url.Values{})
	if errPost != nil {
		zeroLogger.Error().Msgf("Error forking Repo from '%s/%s' ", owner, repo)
		return "", "", errPost
	}
	body, errBody := IoutilReadAll(resp.Body)
	if errBody != nil {
		zeroLogger.Error().Msgf("Error forking Repo from '%s/%s', error %s ", owner, repo, errBody)
		return "", "", errBody
	}
	zeroLogger.Info().Msgf("Status Code after forking the repo: %d", resp.StatusCode)
	formatBody := string(body)
	var result = forkResponseParams{}
	parseError := JSONUnmarshal([]byte(formatBody), &result)
	if parseError != nil {
		zeroLogger.Error().Msgf("Contents seems to be incorrect in the JSON received: %v", parseError)
		return nil, nil, parseError
	}
	zeroLogger.Info().Msgf("User who forked: %s", result.Owner.Login)
	return result.Owner.Login, result.GitURL, err
}

func (gitService gitServiceImplementation) CheckForkedRepo(getURL string) error {
	zeroLogger.Info().Msgf("Checking if repo was forked properly")
	for {
		response, err := HTTPGetCheckForkedRepo(getURL)
		if response.StatusCode == 200 {
			zeroLogger.Info().Msgf("Repo has been forked successfully")
			return nil
		}
		if err != nil || response.StatusCode != 200 {
			zeroLogger.Error().Msgf("Error: %v", err)
			return err
		}
	}
}
