package services

import (
	"github.com/google/go-github/v33/github"
	. "github.com/onsi/ginkgo"
)

type gitServiceInt struct {
	gitClientHandler func() *github.Client
}

var _ = Describe("Services", func() {

})
