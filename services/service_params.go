package services

type forkResponseParams struct {
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	GitURL string `json:"git_url"`
}
