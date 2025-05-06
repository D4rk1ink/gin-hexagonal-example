package hash_test

import (
	"context"
	"testing"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hash Suite")
}

var (
	ctx context.Context

	hashInstance hash.Hash
)

var _ = BeforeSuite(func() {
	ctx = context.TODO()

	hashInstance = hash.NewHash()
})
