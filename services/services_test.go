package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)



type gitServiceMock struct {
	HeadHandler         func(*git.Repository) (*plumbing.Reference, error)
	WorktreeHandler     func(*git.Repository) (*git.Worktree, error)
	FetchHandler        func(repoGit *git.Repository) error
	CheckoutHandler     func(*git.Worktree, string, *plumbing.Reference) (error, bool)
	AddHandler          func(*git.Worktree, string) (plumbing.Hash, error)
	CommitHandler       func(*git.Worktree, string, string) (plumbing.Hash, error)
	CommitObjectHandler func(*git.Repository, plumbing.Hash) (*object.Commit, error)
	PushHandler         func(*git.Repository) error
	PullRequestHandler  func(context.Context, *github.Client, string, string, string, string, string, string, string) (*github.PullRequest, *github.Response, error)
}

// get branch head
func (mock gitServiceMock) Head(repo *git.Repository) (*plumbing.Reference, error) {
	return mock.HeadHandler(repo)
}

// get branch worktree
func (mock gitServiceMock) Worktree(repo *git.Repository) (*git.Worktree, error) {
	return mock.WorktreeHandler(repo)
}

func (mock gitServiceMock) Fetch(repo *git.Repository) error {
	return mock.FetchHandler(repo)
}

func (mock gitServiceMock) Checkout(workingBranch *git.Worktree, branch string, ref *plumbing.Reference) (error,bool) {
	return mock.CheckoutHandler(workingBranch, branch, ref)
}

func (mock gitServiceMock) Add (workingBranch *git.Worktree, secretFile string) (plumbing.Hash, error) {
	return mock.AddHandler(workingBranch, secretFile)
}

func (mock gitServiceMock) Commit(workingBranch *git.Worktree, msg string, owner string) (plumbing.Hash, error) {
	return mock.CommitHandler(workingBranch, msg, owner)
}

func (mock gitServiceMock) CommitObject(repoGit *git.Repository, commit plumbing.Hash) (*object.Commit, error) {
	return mock.CommitObjectHandler(repoGit, commit)
}

func (mock gitServiceMock) Push(repoGit *git.Repository) error {
	return mock.PushHandler(repoGit)
}

func (mock gitServiceMock) PullRequest(ctx context.Context,
	client *github.Client,
	originalOwner string,
	repo string,
	action string,
	owner string,
	currentBranch string,
	headBranch string,
	description string) (*github.PullRequest, *github.Response, error) {
	return mock.PullRequestHandler(
		ctx,
		client,
		originalOwner,
		repo,
		action,
		owner,
		currentBranch,
		headBranch,
		description)
}

var _ = Describe("Service=>", func() {
	Context("GetGitHubClient=>", func() {
		Context(" when trying to create new github client=>", func() {
			It("verifies user by token and returns new github client", func() {
				BackgroundMock := ContextBackground
				defer func() { ContextBackground = BackgroundMock }()
				ContextBackground = func() context.Context {
					var ctx context.Context
					return ctx
				}
				StaticTokenSourceMock := Oauth2StaticTokenSource
				defer func() { Oauth2StaticTokenSource = StaticTokenSourceMock }()
				Oauth2StaticTokenSource = func(*oauth2.Token) oauth2.TokenSource {
					var src oauth2.TokenSource
					return src
				}
				NewClientMock := Oauth2NewClient
				defer func() { Oauth2NewClient = NewClientMock }()
				Oauth2NewClient = func(context.Context, oauth2.TokenSource) *http.Client {
					return new(http.Client)
				}
				GitNewClientMock := GithubNewClient
				defer func() { GithubNewClient = GitNewClientMock }()
				GithubNewClient = func(*http.Client) *github.Client {
					return new(github.Client)
				}
				result := GitServiceObject.GetGitHubClient()
				Expect(result).To(Equal(new(github.Client)))
			})
		})
	})

	Context("CheckUserAccessRepo=>", func() {
		Context("when given user access token=>", func(){
			It("returns repo info", func(){
				GithubRepositoriesMock := GithubRepositories
				defer func() { GithubRepositories = GithubRepositoriesMock }()
				GithubRepositories = func(context.Context, string, string) (*github.Repository, *github.Response, error) {
					return new(github.Repository), new(github.Response), nil
				}
				result, err := GitServiceObject.CheckUserAccessRepo("owner_test", "repo_test")
				Expect(err).To(BeNil())
				Expect(result).To(Equal(new(github.Repository)))
			})
		})

		Context("when user doesn't have rights to access  repo=>", func() {
			It("returns error and exits method", func(){

			})
		})
	})

	Context("CloneRepo=>", func(){
		Context("when owner and repo are given=>", func() {
			It("returns repo info and local directory path", func() {
				GitPlainCloneMock := GitPlainClone
				defer func() { GitPlainClone = GitPlainCloneMock }()
				GitPlainClone = func(string, bool, *git.CloneOptions) (*git.Repository, error) {
					return nil, nil
				}
				_, _, err := GitServiceObject.CloneRepo("owner", "repo")
				Expect(err).To(BeNil())
			})
			Context("if problem occurs in finding the path=>", func() {
				It("returns error and exits the program", func() {
					OsIsNotExistMock := OsIsNotExist
					defer func() { OsIsNotExist = OsIsNotExistMock }()
					OsIsNotExist = func(error) bool {
						return false
					}
					OsRemoveAllMock := OsRemoveAll
					defer func() { OsRemoveAll = OsRemoveAllMock }()
					OsRemoveAll = func(string) error {
						return errors.New("Error in OsRemoveAll")
					}
					_, _, err := GitServiceObject.CloneRepo("owner", "repo")
					Expect(err).To(Equal(errors.New("Error in OsRemoveAll")))
				})
			})
			Context("if problem occurs in cloning the repo=>", func() {
				It("returns error and exits the program", func() {
					GitPlainCloneMock := GitPlainClone
					defer func() { GitPlainClone = GitPlainCloneMock }()
					GitPlainClone = func(string, bool, *git.CloneOptions) (*git.Repository, error) {
						return nil, errors.New("Error in GitPlainClone")
					}
					_, _, err := GitServiceObject.CloneRepo("owner", "repo")
					Expect(err).To(Equal(errors.New("Error in GitPlainClone")))
				})
			})
		})
	})

	// create branch tests
	Context("CreateBranchRepo=>", func() {
		Context("when providing repo details=>", func(){
			It("returns current branch and head branch", func(){
				gitServiceObj := gitServiceMock{}

				gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
					return new(plumbing.Reference), nil
				}
				gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
					return new(git.Worktree), nil
				}
				gitServiceObj.FetchHandler = func(repo *git.Repository) error {
					return nil
				}
				gitServiceObj.CheckoutHandler = func(*git.Worktree, string, *plumbing.Reference) (error, bool) {
					return nil, true
				}
				ThirdPartyGitHub = gitServiceObj
				branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
				Expect(err).To(BeNil())
				Expect(strings.Contains(branch, "secrets_baseline_file")).To(BeTrue())
				Expect(head).To(Equal(""))
			})

			Context("when when problem occurs in fetching branch head=>", func() {
				It("returns error and exits program", func(){
					gitServiceObj := gitServiceMock{}

					gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
						return nil , errors.New("error creating branch")
					}

					ThirdPartyGitHub = gitServiceObj
					branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
					Expect(strings.Contains(fmt.Sprintf("%v",err), "error")).To(BeTrue())
					Expect(branch).To(Equal(""))
					Expect(head).To(Equal(""))
				})
			})

			Context("problem occurs in fetching working branch=>", func() {
				It("returns error and exits program", func(){
					gitServiceObj := gitServiceMock{}

					gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
						return new(plumbing.Reference), nil
					}

					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return nil, errors.New("error fetching working branch")
					}

					ThirdPartyGitHub = gitServiceObj
					branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
					Expect(strings.Contains(fmt.Sprintf("%v",err), "error fetching working branch")).To(BeTrue())
					Expect(branch).To(Equal(""))
					Expect(head).To(Equal(""))
				})
			})

			Context("when problem occurs in fetching all branches from remote for repo=>", func(){
				It("returns error and exits program", func(){
					gitServiceObj := gitServiceMock{}

					gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
						return new(plumbing.Reference), nil
					}

					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}

					gitServiceObj.FetchHandler = func(repo *git.Repository) error {
						return errors.New("error fetching all branches")
					}

					ThirdPartyGitHub = gitServiceObj
					branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
					Expect(strings.Contains(fmt.Sprintf("%v",err), "error fetching all branches")).To(BeTrue())
					Expect(branch).To(Equal(""))
					Expect(head).To(Equal(""))
				})
			})

			Context("when problem occurs in checkout to a branch=>", func() {
				It("returns error and exits program", func(){
					gitServiceObj := gitServiceMock{}

					gitServiceObj.HeadHandler = func(*git.Repository) (*plumbing.Reference, error) {
						return new(plumbing.Reference), nil
					}

					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}

					gitServiceObj.FetchHandler = func(repo *git.Repository) error {
						return nil
					}

					gitServiceObj.CheckoutHandler = func(*git.Worktree, string, *plumbing.Reference) (error, bool) {
						return errors.New("error checkout branch"), false
					}

					ThirdPartyGitHub = gitServiceObj
					branch, head, err:= GitServiceObject.CreateBranchRepo(new(git.Repository), "repo", "create")
					Expect(strings.Contains(fmt.Sprintf("%v",err), "error checkout branch")).To(BeTrue())
					Expect(branch).To(Equal(""))
					Expect(head).To(Equal(""))
				})
			})
		})
	})

	// createSecretFile Test
	Context("CreateSecretFile=>", func(){
		Context("when folder path and secret contents are given=>", func(){
			It("create secrets file at the given path", func(){
				OsStatMock := OsStat
				defer func() { OsStat = OsStatMock }()
				OsStat = func(string) (os.FileInfo, error) {
					return nil, nil
				}
				OsIsNotExistMock := OsIsNotExist
				defer func() { OsIsNotExist = OsIsNotExistMock }()
				OsIsNotExist = func(error) bool {
					return true
				}
				IoutilWriteFileMock := IoutilWriteFile
				defer func() { IoutilWriteFile = IoutilWriteFileMock }()
				IoutilWriteFile = func(filename string, data []byte, perm os.FileMode) error {
					return nil
				}
				err := GitServiceObject.CreateSecretFile("path", "name_file")
				Expect(err).To(BeNil())
			})

			Context("when path doesn't exists=>", func(){
				It("returns error and exits program", func(){
					OsStatMock := OsStat
					defer func() { OsStat = OsStatMock }()
					OsStat = func(string) (os.FileInfo, error) {
						return nil, errors.New("Error in OsStat")
					}
					OsIsNotExistMock := OsIsNotExist
					defer func() { OsIsNotExist = OsIsNotExistMock }()
					OsIsNotExist = func(error) bool {
						return true
					}
					err := GitServiceObject.CreateSecretFile("path", "name_file")
					Expect(err).To(Equal(errors.New("Error in OsStat")))
				})
			})

			Context("when error writing the file=>", func(){
				It("returns error and exits program", func(){
					OsStatMock := OsStat
					defer func() { OsStat = OsStatMock }()
					OsStat = func(string) (os.FileInfo, error) {
						return nil, errors.New("Error in OsStat")
					}
					IoutilWriteFileMock := IoutilWriteFile
					defer func() { IoutilWriteFile = IoutilWriteFileMock }()
					IoutilWriteFile = func(filename string, data []byte, perm os.FileMode) error {
						return errors.New("Error in IoutilWriteFile")
					}
					err := GitServiceObject.CreateSecretFile("path", "name_file")
					Expect(err).To(Equal(errors.New("Error in IoutilWriteFile")))
				})
			})
		})
	})


	// edit secret file Tests
	Context("EditSecretFile=>", func(){
		Context("when folder path and changed contents are given=>", func() {
			It("save changes in the file", func(){
				json1 := `{
                "test.py": [
              {
                "hashed_secret": "b469e216215a36755efee38edf108578280a7f12",
                "is_verified": true,
                        "is_secret": false,
                "line_number": 1,
                "type": "Secret Keyword"
              }
                    ],
                "vars/aws_sdc.yml": [
              {
                "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                "is_verified": true,
                        "is_secret": true,
                "line_number": 9,
                "type": "Secret Keyword"
              }
                    ]
            }`
				json2 := `{
                "custom_plugin_paths": [],
                "exclude": {
                  "files": null,
                  "lines": null
                },
                "generated_at": "2021-01-12T00:00:05Z",
                "plugins_used": [
                  {
                    "name": "AWSKeyDetector"
                  },
                  {
                    "name": "ArtifactoryDetector"
                  },
                  {
                    "base64_limit": 4.5,
                    "name": "Base64HighEntropyString"
                  },
                  {
                    "name": "BasicAuthDetector"
                  },
                  {
                    "name": "CloudantDetector"
                  },
                  {
                    "hex_limit": 3,
                    "name": "HexHighEntropyString"
                  },
                  {
                    "name": "IbmCloudIamDetector"
                  },
                  {
                    "name": "IbmCosHmacDetector"
                  },
                  {
                    "name": "JwtTokenDetector"
                  },
                  {
                    "keyword_exclude": null,
                    "name": "KeywordDetector"
                  },
                  {
                    "name": "MailchimpDetector"
                  },
                  {
                    "name": "PrivateKeyDetector"
                  },
                  {
                    "name": "SlackDetector"
                  },
                  {
                    "name": "SoftlayerDetector"
                  },
                  {
                    "name": "StripeDetector"
                  },
                  {
                    "name": "TwilioKeyDetector"
                  }
                ],
                "results": {
                  "vars/aws_horizontal.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ],
                  "vars/aws_sdc.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ]
                },
                "version": "0.14.3",
                "word_list": {
                  "file": null,
                  "hash": null
                }
              }`
				ioutilReadFileMock := ioutilReadFile
				defer func() { ioutilReadFile = ioutilReadFileMock }()
				ioutilReadFile = func(string) ([]byte, error) {
					return []byte(json2), nil
				}
				JSONMarshalIndentMock := JSONMarshalIndent
				defer func() { JSONMarshalIndent = JSONMarshalIndentMock }()
				JSONMarshalIndent = func(v interface{}, prefix, indent string) ([]byte, error) {
					return []byte(json2), nil
				}
				IoutilWriteFileMock := IoutilWriteFile
				defer func() { IoutilWriteFile = IoutilWriteFileMock }()
				IoutilWriteFile = func(filename string, data []byte, perm os.FileMode) error {
					return nil
				}
				type UpdateParams struct {
					Repo    string                              `json:"repo" xml:"repo" form:"repo"`
					Owner   string                              `json:"owner" xml:"owner" form:"owner"`
					Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
				}
				var secretsChangesMock map[string][]map[string]interface{}
				_ = json.Unmarshal([]byte(json1), &secretsChangesMock)
				err := GitServiceObject.EditSecretFile("owner", secretsChangesMock)
				Expect(err).To(BeNil())
			})

			Context("when file is  not found at the path=>", func(){
				It("returns error and exits the program", func(){
					json1 := `{
                "test.py": [
              {
                "hashed_secret": "b469e216215a36755efee38edf108578280a7f12",
                "is_verified": true,
                        "is_secret": false,
                "line_number": 1,
                "type": "Secret Keyword"
              }
                    ],
                "vars/aws_sdc.yml": [
              {
                "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                "is_verified": true,
                        "is_secret": true,
                "line_number": 9,
                "type": "Secret Keyword"
              }
                    ]
            }`
					json2 := `{
                "custom_plugin_paths": [],
                "exclude": {
                  "files": null,
                  "lines": null
                },
                "generated_at": "2021-01-12T00:00:05Z",
                "plugins_used": [
                  {
                    "name": "AWSKeyDetector"
                  },
                  {
                    "name": "ArtifactoryDetector"
                  },
                  {
                    "base64_limit": 4.5,
                    "name": "Base64HighEntropyString"
                  },
                  {
                    "name": "BasicAuthDetector"
                  },
                  {
                    "name": "CloudantDetector"
                  },
                  {
                    "hex_limit": 3,
                    "name": "HexHighEntropyString"
                  },
                  {
                    "name": "IbmCloudIamDetector"
                  },
                  {
                    "name": "IbmCosHmacDetector"
                  },
                  {
                    "name": "JwtTokenDetector"
                  },
                  {
                    "keyword_exclude": null,
                    "name": "KeywordDetector"
                  },
                  {
                    "name": "MailchimpDetector"
                  },
                  {
                    "name": "PrivateKeyDetector"
                  },
                  {
                    "name": "SlackDetector"
                  },
                  {
                    "name": "SoftlayerDetector"
                  },
                  {
                    "name": "StripeDetector"
                  },
                  {
                    "name": "TwilioKeyDetector"
                  }
                ],
                "results": {
                  "vars/aws_horizontal.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ],
                  "vars/aws_sdc.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ]
                },
                "version": "0.14.3",
                "word_list": {
                  "file": null,
                  "hash": null
                }
              }`
					ioutilReadFileMock := ioutilReadFile
					defer func() { ioutilReadFile = ioutilReadFileMock }()
					ioutilReadFile = func(string) ([]byte, error) {
						return []byte(json2), errors.New("Error in ioutilReadFile")
					}
					type UpdateParams struct {
						Repo    string                              `json:"repo" xml:"repo" form:"repo"`
						Owner   string                              `json:"owner" xml:"owner" form:"owner"`
						Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
					}
					var secretsChangesMock map[string][]map[string]interface{}
					_ = json.Unmarshal([]byte(json1), &secretsChangesMock)
					err := GitServiceObject.EditSecretFile("owner", secretsChangesMock)
					Expect(err).To(Equal(errors.New("Error in ioutilReadFile")))
				})
			})

			Context("when there is problem in getting JSON from existing secrets file", func(){
				It("returns error and exits the program", func(){
					json1 := `{
                "test.py": [
              {
                "hashed_secret": "b469e216215a36755efee38edf108578280a7f12",
                "is_verified": true,
                        "is_secret": false,
                "line_number": 1,
                "type": "Secret Keyword"
              }
                    ],
                "vars/aws_sdc.yml": [
              {
                "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                "is_verified": true,
                        "is_secret": true,
                "line_number": 9,
                "type": "Secret Keyword"
              }
                    ]
            }`
					json2 := `{
                "custom_plugin_paths": [],
                "exclude": {
                  "files": null,
                  "lines": null
                },
                "generated_at": "2021-01-12T00:00:05Z",
                "plugins_used": [
                  {
                    "name": "AWSKeyDetector"
                  },
                  {
                    "name": "ArtifactoryDetector"
                  },
                  {
                    "base64_limit": 4.5,
                    "name": "Base64HighEntropyString"
                  },
                  {
                    "name": "BasicAuthDetector"
                  },
                  {
                    "name": "CloudantDetector"
                  },
                  {
                    "hex_limit": 3,
                    "name": "HexHighEntropyString"
                  },
                  {
                    "name": "IbmCloudIamDetector"
                  },
                  {
                    "name": "IbmCosHmacDetector"
                  },
                  {
                    "name": "JwtTokenDetector"
                  },
                  {
                    "keyword_exclude": null,
                    "name": "KeywordDetector"
                  },
                  {
                    "name": "MailchimpDetector"
                  },
                  {
                    "name": "PrivateKeyDetector"
                  },
                  {
                    "name": "SlackDetector"
                  },
                  {
                    "name": "SoftlayerDetector"
                  },
                  {
                    "name": "StripeDetector"
                  },
                  {
                    "name": "TwilioKeyDetector"
                  }
                ],
                "results": {
                  "vars/aws_horizontal.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ],
                  "vars/aws_sdc.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ]
                },
                "version": "0.14.3",
                "word_list": {
                  "file": null,
                  "hash": null
                }
              }`
					ioutilReadFileMock := ioutilReadFile
					defer func() { ioutilReadFile = ioutilReadFileMock }()
					ioutilReadFile = func(string) ([]byte, error) {
						return []byte(json2), nil
					}
					JSONUnmarshalMock := JSONUnmarshal
					defer func() { JSONUnmarshal = JSONUnmarshalMock }()
					JSONUnmarshal = func([]byte, interface{}) error {
						return errors.New("Error in JSONUnmarshal")
					}
					type UpdateParams struct {
						Repo    string                              `json:"repo" xml:"repo" form:"repo"`
						Owner   string                              `json:"owner" xml:"owner" form:"owner"`
						Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
					}
					var secretsChangesMock map[string][]map[string]interface{}
					_= json.Unmarshal([]byte(json1), &secretsChangesMock)
					err := GitServiceObject.EditSecretFile("owner", secretsChangesMock)
					Expect(err).To(Equal(errors.New("Error in JSONUnmarshal")))
				})
			})

			Context("when there is a problem in parsing the contents of the file=>", func(){
				It("returns error and exits the program", func(){
					json1 := `{
                "test.py": [
              {
                "hashed_secret": "b469e216215a36755efee38edf108578280a7f12",
                "is_verified": true,
                        "is_secret": false,
                "line_number": 1,
                "type": "Secret Keyword"
              }
                    ],
                "vars/aws_sdc.yml": [
              {
                "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                "is_verified": true,
                        "is_secret": true,
                "line_number": 9,
                "type": "Secret Keyword"
              }
                    ]
            }`
					json2 := `{
                "custom_plugin_paths": [],
                "exclude": {
                  "files": null,
                  "lines": null
                },
                "generated_at": "2021-01-12T00:00:05Z",
                "plugins_used": [
                  {
                    "name": "AWSKeyDetector"
                  },
                  {
                    "name": "ArtifactoryDetector"
                  },
                  {
                    "base64_limit": 4.5,
                    "name": "Base64HighEntropyString"
                  },
                  {
                    "name": "BasicAuthDetector"
                  },
                  {
                    "name": "CloudantDetector"
                  },
                  {
                    "hex_limit": 3,
                    "name": "HexHighEntropyString"
                  },
                  {
                    "name": "IbmCloudIamDetector"
                  },
                  {
                    "name": "IbmCosHmacDetector"
                  },
                  {
                    "name": "JwtTokenDetector"
                  },
                  {
                    "keyword_exclude": null,
                    "name": "KeywordDetector"
                  },
                  {
                    "name": "MailchimpDetector"
                  },
                  {
                    "name": "PrivateKeyDetector"
                  },
                  {
                    "name": "SlackDetector"
                  },
                  {
                    "name": "SoftlayerDetector"
                  },
                  {
                    "name": "StripeDetector"
                  },
                  {
                    "name": "TwilioKeyDetector"
                  }
                ],
                "version": "0.14.3",
                "word_list": {
                  "file": null,
                  "hash": null
                }
              }`
					ioutilReadFileMock := ioutilReadFile
					defer func() { ioutilReadFile = ioutilReadFileMock }()
					ioutilReadFile = func(string) ([]byte, error) {
						return []byte(json2), nil
					}
					//type UpdateParams struct {
					//	Repo    string                              `json:"repo" xml:"repo" form:"repo"`
					//	Owner   string                              `json:"owner" xml:"owner" form:"owner"`
					//	Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
					//}
					var secretsChangesMock map[string][]map[string]interface{}
					_ = json.Unmarshal([]byte(json1), &secretsChangesMock)
					err := GitServiceObject.EditSecretFile("owner", secretsChangesMock)
					Expect(err).To(Equal(errors.New("could not parse the result data in secret file, please check the data")))
				})
			})

			Context("when indentation fails  while editing the file=>", func(){
				It("returns error and exits the program", func(){
					json1 := `{
                "test.py": [
              {
                "hashed_secret": "b469e216215a36755efee38edf108578280a7f12",
                "is_verified": true,
                        "is_secret": false,
                "line_number": 1,
                "type": "Secret Keyword"
              }
                    ],
                "vars/aws_sdc.yml": [
              {
                "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                "is_verified": true,
                        "is_secret": true,
                "line_number": 9,
                "type": "Secret Keyword"
              }
                    ]
            }`
					json2 := `{
                "custom_plugin_paths": [],
                "exclude": {
                  "files": null,
                  "lines": null
                },
                "generated_at": "2021-01-12T00:00:05Z",
                "plugins_used": [
                  {
                    "name": "AWSKeyDetector"
                  },
                  {
                    "name": "ArtifactoryDetector"
                  },
                  {
                    "base64_limit": 4.5,
                    "name": "Base64HighEntropyString"
                  },
                  {
                    "name": "BasicAuthDetector"
                  },
                  {
                    "name": "CloudantDetector"
                  },
                  {
                    "hex_limit": 3,
                    "name": "HexHighEntropyString"
                  },
                  {
                    "name": "IbmCloudIamDetector"
                  },
                  {
                    "name": "IbmCosHmacDetector"
                  },
                  {
                    "name": "JwtTokenDetector"
                  },
                  {
                    "keyword_exclude": null,
                    "name": "KeywordDetector"
                  },
                  {
                    "name": "MailchimpDetector"
                  },
                  {
                    "name": "PrivateKeyDetector"
                  },
                  {
                    "name": "SlackDetector"
                  },
                  {
                    "name": "SoftlayerDetector"
                  },
                  {
                    "name": "StripeDetector"
                  },
                  {
                    "name": "TwilioKeyDetector"
                  }
                ],
                "results": {
                  "vars/aws_horizontal.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ],
                  "vars/aws_sdc.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ]
                },
                "version": "0.14.3",
                "word_list": {
                  "file": null,
                  "hash": null
                }
              }`
					ioutilReadFileMock := ioutilReadFile
					defer func() { ioutilReadFile = ioutilReadFileMock }()
					ioutilReadFile = func(string) ([]byte, error) {
						return []byte(json2), nil
					}
					JSONMarshalIndentMock := JSONMarshalIndent
					defer func() { JSONMarshalIndent = JSONMarshalIndentMock }()
					JSONMarshalIndent = func(v interface{}, prefix, ident string) ([]byte, error) {
						return []byte(json2), errors.New("Error in JSONMarshalIndent")
					}
					type UpdateParams struct {
						Repo    string                              `json:"repo" xml:"repo" form:"repo"`
						Owner   string                              `json:"owner" xml:"owner" form:"owner"`
						Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
					}
					var secretsChangesMock map[string][]map[string]interface{}
					_ = json.Unmarshal([]byte(json1), &secretsChangesMock)
					err := GitServiceObject.EditSecretFile("owner", secretsChangesMock)
					Expect(err).To(Equal(errors.New("Error in JSONMarshalIndent")))
				})
			})

			Context("when problem arises while  writing file=>", func(){
				It("returns error and exits the program", func(){
					json1 := `{
                "test.py": [
              {
                "hashed_secret": "b469e216215a36755efee38edf108578280a7f12",
                "is_verified": true,
                        "is_secret": false,
                "line_number": 1,
                "type": "Secret Keyword"
              }
                    ],
                "vars/aws_sdc.yml": [
              {
                "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                "is_verified": true,
                        "is_secret": true,
                "line_number": 9,
                "type": "Secret Keyword"
              }
                    ]
            }`
					json2 := `{
                "custom_plugin_paths": [],
                "exclude": {
                  "files": null,
                  "lines": null
                },
                "generated_at": "2021-01-12T00:00:05Z",
                "plugins_used": [
                  {
                    "name": "AWSKeyDetector"
                  },
                  {
                    "name": "ArtifactoryDetector"
                  },
                  {
                    "base64_limit": 4.5,
                    "name": "Base64HighEntropyString"
                  },
                  {
                    "name": "BasicAuthDetector"
                  },
                  {
                    "name": "CloudantDetector"
                  },
                  {
                    "hex_limit": 3,
                    "name": "HexHighEntropyString"
                  },
                  {
                    "name": "IbmCloudIamDetector"
                  },
                  {
                    "name": "IbmCosHmacDetector"
                  },
                  {
                    "name": "JwtTokenDetector"
                  },
                  {
                    "keyword_exclude": null,
                    "name": "KeywordDetector"
                  },
                  {
                    "name": "MailchimpDetector"
                  },
                  {
                    "name": "PrivateKeyDetector"
                  },
                  {
                    "name": "SlackDetector"
                  },
                  {
                    "name": "SoftlayerDetector"
                  },
                  {
                    "name": "StripeDetector"
                  },
                  {
                    "name": "TwilioKeyDetector"
                  }
                ],
                "results": {
                  "vars/aws_horizontal.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ],
                  "vars/aws_sdc.yml": [
                    {
                      "hashed_secret": "9b1a33f9e0bbed1e3def05152c38d9312611f166",
                      "is_verified": false,
                      "line_number": 9,
                      "type": "Secret Keyword"
                    }
                  ]
                },
                "version": "0.14.3",
                "word_list": {
                  "file": null,
                  "hash": null
                }
              }`
					ioutilReadFileMock := ioutilReadFile
					defer func() { ioutilReadFile = ioutilReadFileMock }()
					ioutilReadFile = func(string) ([]byte, error) {
						return []byte(json2), nil
					}
					JSONMarshalIndentMock := JSONMarshalIndent
					defer func() { JSONMarshalIndent = JSONMarshalIndentMock }()
					JSONMarshalIndent = func(v interface{}, prefix, ident string) ([]byte, error) {
						return []byte(json2), nil
					}
					IoutilWriteFiletMock := IoutilWriteFile
					defer func() { IoutilWriteFile = IoutilWriteFiletMock }()
					IoutilWriteFile = func(filename string, data []byte, perm os.FileMode) error {
						return errors.New("Error in IoutilWriteFile")
					}
					type UpdateParams struct {
						Repo    string                              `json:"repo" xml:"repo" form:"repo"`
						Owner   string                              `json:"owner" xml:"owner" form:"owner"`
						Changes map[string][]map[string]interface{} `json:"changes" xml:"changes" form:"changes"`
					}
					var secretsChangesMock map[string][]map[string]interface{}
					_ = json.Unmarshal([]byte(json1), &secretsChangesMock)
					err := GitServiceObject.EditSecretFile("owner", secretsChangesMock)
					Expect(err).To(Equal(errors.New("Error in IoutilWriteFile")))
				})
			})
		})
	})

	// CreateCommitAndPr tests
	Context("CreateCommitAndPr=>", func() {
		Context("when provided repo, branch and owner info=>", func() {
			It("commits and push changes, and creates PR", func() {
				gitServiceObj := gitServiceMock{}
				BackgroundMock := ContextBackground
				defer func() { ContextBackground = BackgroundMock }()
				ContextBackground = func() context.Context {
					var ctx context.Context
					return ctx
				}
				StaticTokenSourceMock := Oauth2StaticTokenSource
				defer func() { Oauth2StaticTokenSource = StaticTokenSourceMock }()
				Oauth2StaticTokenSource = func(*oauth2.Token) oauth2.TokenSource {
					var src oauth2.TokenSource
					return src
				}
				NewClientMock := Oauth2NewClient
				defer func() { Oauth2NewClient = NewClientMock }()
				Oauth2NewClient = func(context.Context, oauth2.TokenSource) *http.Client {
					return new(http.Client)
				}
				GitNewClientMock := GithubNewClient
				defer func() { GithubNewClient = GitNewClientMock }()
				GithubNewClient = func(*http.Client) *github.Client {
					return new(github.Client)
				}
				gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
					return new(git.Worktree), nil
				}
				gitServiceObj.AddHandler = func(*git.Worktree, string) (plumbing.Hash, error) {
					return plumbing.NewHash("secret"), nil
				}
				gitServiceObj.CommitHandler = func(*git.Worktree, string, string) (plumbing.Hash, error) {
					return plumbing.NewHash("secret"), nil
				}
				gitServiceObj.CommitObjectHandler = func(*git.Repository, plumbing.Hash) (*object.Commit, error) {
					return new(object.Commit), nil
				}
				gitServiceObj.PushHandler = func(*git.Repository) error {
					return nil
				}
				gitServiceObj.PullRequestHandler = func(context.Context, *github.Client, string, string, string, string, string, string, string) (*github.PullRequest, *github.Response, error) {
					return new(github.PullRequest), new(github.Response), nil
				}

				ThirdPartyGitHub = gitServiceObj

				err := GitServiceObject.CreateCommitAndPr(
					"owner",
					"originalOwner",
					"repo",
					"currentBranch",
					"headBranch",
					"action",
					"description",
					new(git.Repository))
				Expect(err).To(BeNil())
			})

			Context("if problem occurs in fetching working branch=>", func(){
				It("returns error and exists program ", func() {
					gitServiceObj := gitServiceMock{}
					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return nil, errors.New("error  finding working branch")
					}
					ThirdPartyGitHub = gitServiceObj

					err := GitServiceObject.CreateCommitAndPr(
						"owner",
						"originalOwner",
						"repo",
						"currentBranch",
						"headBranch",
						"action",
						"description",
						new(git.Repository))
					Expect(strings.Contains(fmt.Sprintf("%v", err), "error")).To(BeTrue())
				})
			})


			Context("if problem occurs in adding changes to branch=>", func(){
				It("returns error and exists program", func() {
					gitServiceObj := gitServiceMock{}
					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}
					gitServiceObj.AddHandler = func(*git.Worktree, string) (plumbing.Hash, error) {
						return plumbing.NewHash(""), errors.New("error  adding changes")
					}
					ThirdPartyGitHub = gitServiceObj

					err := GitServiceObject.CreateCommitAndPr(
						"owner",
						"originalOwner",
						"repo",
						"currentBranch",
						"headBranch",
						"action",
						"description",
						new(git.Repository))
					Expect(strings.Contains(fmt.Sprintf("%v", err), "error")).To(BeTrue())
				})
			})

			Context("if problem occurs in committing changes=>", func(){
				It("returns error and exists program", func() {
					gitServiceObj := gitServiceMock{}
					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}
					gitServiceObj.AddHandler = func(*git.Worktree, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitHandler = func(*git.Worktree, string, string) (plumbing.Hash, error) {
						return plumbing.NewHash(""), errors.New("error  committing changes")
					}
					ThirdPartyGitHub = gitServiceObj

					err := GitServiceObject.CreateCommitAndPr(
						"owner",
						"originalOwner",
						"repo",
						"currentBranch",
						"headBranch",
						"action",
						"description",
						new(git.Repository))
					Expect(strings.Contains(fmt.Sprintf("%v", err), "error")).To(BeTrue())
				})
			})

			Context("if problem occurs with CommitObject=>", func() {
				It("returns error and exists program", func() {
					gitServiceObj := gitServiceMock{}
					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}
					gitServiceObj.AddHandler = func(*git.Worktree, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitHandler = func(*git.Worktree, string, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitObjectHandler = func(*git.Repository, plumbing.Hash) (*object.Commit, error) {
						return nil, errors.New("error with commit object")
					}

					ThirdPartyGitHub = gitServiceObj

					err := GitServiceObject.CreateCommitAndPr(
						"owner",
						"originalOwner",
						"repo",
						"currentBranch",
						"headBranch",
						"action",
						"description",
						new(git.Repository))
					Expect(strings.Contains(fmt.Sprintf("%v", err), "error")).To(BeTrue())
				})
			})

			Context("if problem occurs in pushing changes=>", func(){
				It("returns error and exists program", func() {
					gitServiceObj := gitServiceMock{}
					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}
					gitServiceObj.AddHandler = func(*git.Worktree, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitHandler = func(*git.Worktree, string, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitObjectHandler = func(*git.Repository, plumbing.Hash) (*object.Commit, error) {
						return new(object.Commit), nil
					}
					gitServiceObj.PushHandler = func(*git.Repository) error {
						return errors.New("error pushing changes")
					}

					ThirdPartyGitHub = gitServiceObj

					err := GitServiceObject.CreateCommitAndPr(
						"owner",
						"originalOwner",
						"repo",
						"currentBranch",
						"headBranch",
						"action",
						"description",
						new(git.Repository))
					Expect(strings.Contains(fmt.Sprintf("%v", err), "error")).To(BeTrue())
				})
			})

			Context("if problem occurs in creating PR=>", func(){
				It("returns error and exits program", func() {
					gitServiceObj := gitServiceMock{}
					BackgroundMock := ContextBackground
					defer func() { ContextBackground = BackgroundMock }()
					ContextBackground = func() context.Context {
						var ctx context.Context
						return ctx
					}
					StaticTokenSourceMock := Oauth2StaticTokenSource
					defer func() { Oauth2StaticTokenSource = StaticTokenSourceMock }()
					Oauth2StaticTokenSource = func(*oauth2.Token) oauth2.TokenSource {
						var src oauth2.TokenSource
						return src
					}
					NewClientMock := Oauth2NewClient
					defer func() { Oauth2NewClient = NewClientMock }()
					Oauth2NewClient = func(context.Context, oauth2.TokenSource) *http.Client {
						return new(http.Client)
					}
					GitNewClientMock := GithubNewClient
					defer func() { GithubNewClient = GitNewClientMock }()
					GithubNewClient = func(*http.Client) *github.Client {
						return new(github.Client)
					}
					gitServiceObj.WorktreeHandler = func(*git.Repository) (*git.Worktree, error) {
						return new(git.Worktree), nil
					}
					gitServiceObj.AddHandler = func(*git.Worktree, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitHandler = func(*git.Worktree, string, string) (plumbing.Hash, error) {
						return plumbing.NewHash("secret"), nil
					}
					gitServiceObj.CommitObjectHandler = func(*git.Repository, plumbing.Hash) (*object.Commit, error) {
						return new(object.Commit), nil
					}
					gitServiceObj.PushHandler = func(*git.Repository) error {
						return nil
					}
					gitServiceObj.PullRequestHandler = func(context.Context, *github.Client, string, string, string, string, string, string, string) (*github.PullRequest, *github.Response, error) {
						return nil, nil, errors.New("PR Updated")
					}

					ThirdPartyGitHub = gitServiceObj

					err := GitServiceObject.CreateCommitAndPr(
						"owner",
						"originalOwner",
						"repo",
						"currentBranch",
						"headBranch",
						"action",
						"description",
						new(git.Repository))
					Expect(strings.Contains(fmt.Sprintf("%v", err), "PR Updated")).To(BeTrue())
				})
			})
		})
	})

	// Fork repo tests
	Context("ForkRepo=>", func(){
		Context("when given owner and repo info=>", func(){
			It("forks the repo and returns owner and repo URL", func(){
				HTTPPostForkRepoMock := HTTPPostForkRepo
				defer func() { HTTPPostForkRepo = HTTPPostForkRepoMock }()
				HTTPPostForkRepo = func(string, url.Values) (resp *http.Response, err error) {
					var a http.Response
					a.StatusCode = 200
					return &a, nil
				}
				IoutilReadAllMock := IoutilReadAll
				defer func() { IoutilReadAll = IoutilReadAllMock }()
				IoutilReadAll = func(io.Reader) ([]byte, error) {
					var t []byte
					return t, nil
				}
				JSONUnmarshalMock := JSONUnmarshal
				defer func() { JSONUnmarshal = JSONUnmarshalMock }()
				JSONUnmarshal = func([]byte, interface{}) error {
					return nil
				}
				_, _, err := GitServiceObject.ForkRepo("test", "test")
				Expect(err).To(BeNil())
			})

			Context("when there is problem in forking the repo at server=>", func(){
				It("returns 400 status code and exits the program", func(){
					HTTPPostForkRepoMock := HTTPPostForkRepo
					defer func() { HTTPPostForkRepo = HTTPPostForkRepoMock }()
					HTTPPostForkRepo = func(string, url.Values) (resp *http.Response, err error) {
						var a http.Response
						a.StatusCode = 400
						return &a, errors.New("Error Forking repo")
					}
					forkedUser, url, err := GitServiceObject.ForkRepo("test", "test")
					Expect(forkedUser).To(Equal(""))
					Expect(url).To(Equal(""))
					Expect(err).To(Equal(errors.New("Error Forking repo")))
				})
			})

			Context("when response from server is not parsed properly=>", func(){
				It("returns error  and exits the program", func(){
					HTTPPostForkRepoMock := HTTPPostForkRepo
					defer func() { HTTPPostForkRepo = HTTPPostForkRepoMock }()
					HTTPPostForkRepo = func(string, url.Values) (resp *http.Response, err error) {
						var a http.Response
						a.StatusCode = 200
						return &a, nil
					}
					IoutilReadAllMock := IoutilReadAll
					defer func() { IoutilReadAll = IoutilReadAllMock }()
					IoutilReadAll = func(io.Reader) ([]byte, error) {
						var t []byte
						return t, errors.New("Error in IoutilReadAll")
					}
					forkedUser, url, err := GitServiceObject.ForkRepo("test", "test")
					Expect(forkedUser).To(Equal(""))
					Expect(url).To(Equal(""))
					Expect(err).To(Equal(errors.New("Error in IoutilReadAll")))
				})
			})

			Context("when the contents are missing from received JSON=>", func(){
				It("returns error and exits program", func(){
					HTTPPostForkRepoMock := HTTPPostForkRepo
					defer func() { HTTPPostForkRepo = HTTPPostForkRepoMock }()
					HTTPPostForkRepo = func(string, url.Values) (resp *http.Response, err error) {
						var a http.Response
						a.StatusCode = 200
						return &a, nil
					}
					IoutilReadAllMock := IoutilReadAll
					defer func() { IoutilReadAll = IoutilReadAllMock }()
					IoutilReadAll = func(io.Reader) ([]byte, error) {
						var t []byte
						return t, nil
					}
					JSONUnmarshalMock := JSONUnmarshal
					defer func() { JSONUnmarshal = JSONUnmarshalMock }()
					JSONUnmarshal = func([]byte, interface{}) error {
						return errors.New("Error in JSONUnmarshal")
					}
					forkedUser, url, err := GitServiceObject.ForkRepo("test", "test")
					Expect(forkedUser).To(BeNil())
					Expect(url).To(BeNil())
					Expect(err).To(Equal(errors.New("Error in JSONUnmarshal")))
				})
			})
		})
	})

	// CheckForkRepo Tests
	Context("CheckForkedRepo=>", func(){
		Context("when repo is forked properly=>", func(){
			It("return 200 status code", func(){
				HTTPGetRepoMock := HTTPGetCheckForkedRepo
				defer func() { HTTPGetCheckForkedRepo = HTTPGetRepoMock }()
				HTTPGetCheckForkedRepo = func(string) (*http.Response, error) {
					var a http.Response
					a.StatusCode = 200
					return &a, nil
				}
				err := GitServiceObject.CheckForkedRepo("test")
				Expect(err).To(BeNil())
			})

			Context("when 400 status codes is received from server=>", func(){
				It("returns that fork is not done for the repo", func(){
					HTTPGetRepoMock := HTTPGetCheckForkedRepo
					defer func() { HTTPGetCheckForkedRepo = HTTPGetRepoMock }()
					HTTPGetCheckForkedRepo = func(string) (*http.Response, error) {
						var a http.Response
						a.StatusCode = 400
						return &a, errors.New("Error checking Forked repo")
					}
					err := GitServiceObject.CheckForkedRepo("test")
					Expect(err).To(Equal(errors.New("Error checking Forked repo")))
				})
			})
		})
	})
})