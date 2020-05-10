package sql_test

import (
	"testing"

	"github.com/novln/macchiato"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestSQL(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	macchiato.RunSpecs(t, "migrations/sql Test Suite")
}
