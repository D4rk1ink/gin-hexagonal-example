package auth_service_test

import (
	"testing"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = BeforeSuite(func() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
})
