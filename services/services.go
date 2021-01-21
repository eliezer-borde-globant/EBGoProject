package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	. "github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// Background: Third Party packages.
var (
	ContextBackground       = context.Background
	Oauth2StaticTokenSource = oauth2.StaticTokenSource
	Oauth2NewClient         = oauth2.NewClient
	GithubNewClient         = github.NewClient
	GithubRepositories      = GitServiceObject.GetGitHubClient().Repositories.Get
)

func (gitService gitServiceImplementation) GetGitHubClient() *github.Client {
	ctx := ContextBackground()
	ts := Oauth2StaticTokenSource(
		&oauth2.Token{AccessToken: GitHubToken},
	)
	tc := Oauth2NewClient(ctx, ts)
	return GithubNewClient(tc)
}

func (gitService gitServiceImplementation) CheckUserAccessRepo(owner string, repo string) (*github.Repository, error) {
	ZeroLogger.Info().Msgf("check user has access to %s/%s", owner, repo)
	ctx := context.Background()
	repoInfo, _, err := GithubRepositories(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	return repoInfo, nil
}

func (gitService gitServiceImplementation) CloneRepo(owner string, repo string) (*git.Repository, string, error) {
	path := fmt.Sprintf("/tmp/%s-%s", owner, repo)

	ZeroLogger.Info().Msgf("Creating folder to clone %s", path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.RemoveAll(path)
		if err != nil {
			ZeroLogger.Error().Msgf("Error path to clone repo from %s/%s, error: %v", owner, repo, err)
			return nil, "", err
		}
	}
	ZeroLogger.Info().Msg("Starting to clone Repo")
	repoInfo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://%s@github.com/%s/%s", GitHubToken, owner, repo),
		Progress: os.Stdout,
	})
	if err != nil {
		ZeroLogger.Error().Msgf("Error Cloning repo from %s/%s, error: %v", owner, repo, err)
		return nil, "", err
	}
	ZeroLogger.Info().Msgf("Repo was cloned")
	return repoInfo, path, nil
}

func (gitService gitServiceImplementation) CreateBranchRepo(repoGit *git.Repository, repoName string, action string) (string, string, error) {
	ZeroLogger.Info().Msgf("Creating Branch to update secret file in repo %s", repoName)
	headRef, err := repoGit.Head()
	if err != nil {
		ZeroLogger.Error().Msgf("Error Creating Branch to update secret file in repo %s, error: %v", repoName, err)
		return "", "", err
	}
	headBranchName := strings.ReplaceAll(headRef.Name().String(), "refs/heads/", "")
	branch := fmt.Sprintf("secret_scanner_api/%s/%s/secrets_baseline_file", repoName, action)
	workingBranch, err := repoGit.Worktree()
	if err != nil {
		ZeroLogger.Error().Msgf("Error Creating Branch to update secret file in repo %s, error: %v", repoName, err)
		return "", "", err
	}
	ZeroLogger.Info().Msgf("Fetching all Branches from %s", repoName)
	err = repoGit.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		ZeroLogger.Error().Msgf("Error fetching remote Branches from repo %s, error: %v", repoName, err)
		return "", "", err
	}
	ZeroLogger.Info().Msgf("Checking if the branch %s exists in %s", branch, repoName)
	err = workingBranch.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err == nil {
		ZeroLogger.Info().Msgf("Branch %s already exists in %s, Checking out...", branch, repoName)
	} else {
		ZeroLogger.Info().Msgf("Creating new branch %s in %s", branch, repoName)
		err = workingBranch.Checkout(&git.CheckoutOptions{
			Hash:   headRef.Hash(),
			Branch: plumbing.NewBranchReferenceName(branch),
			Create: true,
		})
		if err != nil {
			ZeroLogger.Error().Msgf("Error Creating Branch to update secret file in repo %s, error: %v", repoName, err)
			return "", "", err
		}
		ZeroLogger.Info().Msgf("Branch created in (%s) with the name (%s)", repoName, branch)
	}
	return branch, headBranchName, err
}

func (gitService gitServiceImplementation) CreateSecretFile(path string, secretFile string) error {
	ZeroLogger.Info().Msg(fmt.Sprintf("Creating Path %s to add %s file ", path, SecretsFileName))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ZeroLogger.Error().Msgf("Error Creating Path %s, error: %v", path, err)
		return err
	}
	path = fmt.Sprintf("%s/%s", path, SecretsFileName)
	err := ioutil.WriteFile(path, []byte(secretFile), 0644)
	ZeroLogger.Info().Msgf("File was created with the content at path: '%s'", path)
	if err != nil {
		ZeroLogger.Error().Msgf("Error creating %s file in, %s error: %v", SecretsFileName, path, err)
		return err
	}
	return nil
}

func (gitService gitServiceImplementation) EditSecretFile(path string, secretsChanges SecretUpdateMap) error {
	ZeroLogger.Info().Msgf("Starting to edit the secret file at path: '%s'", path)
	path = fmt.Sprintf("%s/%s", path, SecretsFileName)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var fileStruct map[string]interface{}
	err = json.Unmarshal(dat, &fileStruct)
	if err != nil {
		return err
	}
	results, ok := fileStruct["results"].(map[string]interface{})
	if !ok {
		err := errors.New("could not parse the result data in secret file, please check the data")
		ZeroLogger.Error().Msgf("Error: %v", err)
		return err
	}
	for filename, secretData := range secretsChanges {
		_, ok := results[filename]
		if !ok {
			continue
		}
		for _, secret := range secretData {
			fileData := results[filename]
			fileSecrets := reflect.ValueOf(fileData)
			for i := 0; i < fileSecrets.Len(); i++ {
				value := fileSecrets.Index(i)
				secrets := value.Interface().(map[string]interface{})
				if secret["hashed_secret"] == secrets["hashed_secret"] && secret["line_number"] == secrets["line_number"] {
					secrets["is_secret"] = secret["is_secret"]
					value.Set(reflect.ValueOf(secrets))
				}
			}
			results[filename] = fileSecrets.Interface()
		}
	}
	fileStruct["results"] = results
	file, parseError := json.MarshalIndent(fileStruct, "", "  ")
	if parseError != nil {
		ZeroLogger.Error().Msgf("Cannot indent content of the file : %v", parseError)
		return parseError
	}
	writeFileError := ioutil.WriteFile(path, file, 0644)
	if writeFileError != nil {
		ZeroLogger.Error().Msgf("Error writing file: %v", parseError)
		return parseError
	}
	return nil
}

func (gitService gitServiceImplementation) CreateCommitAndPr(owner string, originalOwner string, repo string, currentBranch string, headBranch string, action string, description string, repoGit *git.Repository) error {
	ZeroLogger.Info().Msg("Getting current branch")
	ZeroLogger.Info().Msgf("Head Branch: %s", headBranch)
	ZeroLogger.Info().Msgf("Current Branch: %s", currentBranch)
	workingBranch, err := repoGit.Worktree()
	if err != nil {
		ZeroLogger.Info().Msgf("Error getting current branch '%s/%s'", owner, repo)
		return err
	}
	ZeroLogger.Info().Msgf("Adding %s file to new branch ", SecretsFileName)
	_, err = workingBranch.Add(SecretsFileName)
	if err != nil {
		return err
	}
	ZeroLogger.Info().Msgf("%s was added to stage ", SecretsFileName)
	ZeroLogger.Info().Msg("Committing Changes")
	commit, err := workingBranch.Commit(fmt.Sprintf("chore: %s secret baseline file", action), &git.CommitOptions{
		Author: &object.Signature{
			Name: owner,
			When: time.Now(),
		},
	})
	if err != nil {
		ZeroLogger.Error().Msgf("Error Committing changes: %v", err)
		return err
	}
	ZeroLogger.Info().Msg("Changes were committed")
	_, err = repoGit.CommitObject(commit)
	if err != nil {
		ZeroLogger.Error().Msgf("Error Committing: %v", err)
		return err
	}
	ZeroLogger.Info().Msgf("Commit created in '%s/%s'", owner, repo)

	ZeroLogger.Info().Msg("Pushing changes to remote")
	obj := &git.PushOptions{}
	err = repoGit.Push(obj)
	if err != nil {
		return err
	}
	ZeroLogger.Info().Msgf("Branch was pushed '%s/%s'", owner, repo)
	githubClient := GitServiceObject.GetGitHubClient()
	newPR := &github.NewPullRequest{
		Title:               github.String(fmt.Sprintf("[Detect Secrets] %s Secret BaseLine File", action)),
		Head:                github.String(fmt.Sprintf("%s:%s", owner, currentBranch)),
		Base:                github.String(headBranch),
		Body:                github.String(description),
		MaintainerCanModify: github.Bool(true),
	}

	_, _, err = githubClient.PullRequests.Create(context.Background(), originalOwner, repo, newPR)
	if err != nil {
		ZeroLogger.Info().Msgf("PR success Updated! '%s/%s'", owner, repo)
		return err
	}
	ZeroLogger.Info().Msgf("PR success Created! '%s/%s'", owner, repo)
	return nil
}

func (gitService gitServiceImplementation) ForkRepo(owner string, repo string) (forkedOwner interface{}, gitURL interface{}, err error) {
	ZeroLogger.Info().Msgf("Forking repo from '%s/%s'", owner, repo)
	forkURL := fmt.Sprintf("https://%s@api.github.com/repos/%s/%s/forks", GitHubToken, owner, repo)
	resp, errPost := http.PostForm(forkURL, url.Values{})

	if errPost != nil {
		ZeroLogger.Error().Msgf("Error forking Repo from '%s/%s' ", owner, repo)
		return "", "", errPost
	}

	body, errBody := ioutil.ReadAll(resp.Body)

	if errBody != nil {
		ZeroLogger.Error().Msgf("Error forking Repo from '%s/%s', error %s ", owner, repo, errBody)
		return "", "", errBody
	}

	ZeroLogger.Info().Msgf("Status Code after forking the repo: %d", resp.StatusCode)
	formatBody := string(body)

	var result = forkResponseParams{}

	parseError := json.Unmarshal([]byte(formatBody), &result)
	if parseError != nil {
		ZeroLogger.Error().Msgf("Contents seems to be incorrect in the JSON received: %v", parseError)
		return nil, nil, parseError
	}
	ZeroLogger.Info().Msgf("User who forked: %s", result.Owner.Login)

	return result.Owner.Login, result.GitURL, err
}

func (gitService gitServiceImplementation) CheckForkedRepo(getURL string) error {
	ZeroLogger.Info().Msgf("Checking if repo was forked properly")
	for {
		response, err := http.Get(getURL)
		if response.StatusCode == 200 {
			ZeroLogger.Info().Msgf("Repo has been forked successfully")
			break
		}
		if err != nil {
			ZeroLogger.Error().Msgf("Error: %v", err)
			return err
		}
	}
	return nil
}
