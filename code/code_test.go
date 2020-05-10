package code_test

import (
	"testing"

	"github.com/novln/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestCode(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "migrations/code Test Suite")
}
