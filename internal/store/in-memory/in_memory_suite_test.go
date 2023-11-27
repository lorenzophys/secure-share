package memory_store_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "InMemory Suite")
}
