package jwt_test

import (
	"testing"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jwt Suite")
}

var jwtInstance jwt.Jwt

var _ = BeforeSuite(func() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	jwtInstance = jwt.NewJwt()
})
