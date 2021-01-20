package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoProject Suite")
}
