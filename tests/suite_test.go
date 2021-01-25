package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEBGoProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EBGoProject Suite")
}
