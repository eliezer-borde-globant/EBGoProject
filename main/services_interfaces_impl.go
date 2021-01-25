package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v33/github"
)

type gitServiceImplementation struct{}
type thirdPartyGitHubImpl struct{}

//GitServiceObject : Object with service methods.
var (
	GitServiceObject gitServiceInterface       = gitServiceImplementation{}
	ThirdPartyGitHub thirdPartyGitHubInterface = thirdPartyGitHubImpl{}
)

func (service thirdPartyGitHubImpl) Head(repoGit *git.Repository) (*plumbing.Reference, error) {
	return repoGit.Head()
}

func (service thirdPartyGitHubImpl) Worktree(repoGit *git.Repository) (*git.Worktree, error) {
	return repoGit.Worktree()
}

func (service thirdPartyGitHubImpl) Fetch(repoGit *git.Repository) error {
	return repoGit.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
}

func (service thirdPartyGitHubImpl) Checkout(workingBranch *git.Worktree, branch string, headRef *plumbing.Reference) (string, error) {
	var branchStatus string
	err := workingBranch.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err == nil {
		branchStatus = "EXISTING_BRANCH"
	} else {
		zeroLogger.Info().Msgf("Creating new branch")

		err = workingBranch.Checkout(&git.CheckoutOptions{
			Hash:   headRef.Hash(),
			Branch: plumbing.NewBranchReferenceName(branch),
			Create: true,
		})
		if err != nil {
			zeroLogger.Fatal().Msgf("Error Creating Branch to update secret file, error: %v", err)
			branchStatus = "ERROR"
		}
		branchStatus = "NEW_BRANCH"
		zeroLogger.Info().Msg("Branch created!")
	}
	return branchStatus, err
}

func (service thirdPartyGitHubImpl) Add(workingBranch *git.Worktree, SecretFile string) (plumbing.Hash, error) {
	return workingBranch.Add(SecretFile)
}

func (service thirdPartyGitHubImpl) Commit(workingBranch *git.Worktree, msg string, owner string) (plumbing.Hash, error) {
	return workingBranch.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name: owner,
			When: time.Now(),
		},
	})
}

func (service thirdPartyGitHubImpl) CommitObject(repoGit *git.Repository, commit plumbing.Hash) (*object.Commit, error) {
	return repoGit.CommitObject(commit)
}

func (service thirdPartyGitHubImpl) Push(repoGit *git.Repository) error {
	obj := &git.PushOptions{}
	return repoGit.Push(obj)
}

func (service thirdPartyGitHubImpl) PullRequest(
	ctx context.Context,
	client *github.Client,
	originalOwner string,
	repo string,
	action string,
	owner string,
	currentBranch string,
	headBranch string,
	description string) (*github.PullRequest, *github.Response, error) {
	newPR := &github.NewPullRequest{
		Title:               github.String(fmt.Sprintf("[Detect Secrets] %s Secret BaseLine File", action)),
		Head:                github.String(fmt.Sprintf("%s:%s", owner, currentBranch)),
		Base:                github.String(headBranch),
		Body:                github.String(description),
		MaintainerCanModify: github.Bool(true),
	}

	return client.PullRequests.Create(ctx, originalOwner, repo, newPR)
}
