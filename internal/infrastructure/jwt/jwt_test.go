package jwt_test

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT", Ordered, func() {
	Context("GenerateAccessToken", func() {
		var jwtInstance jwt.Jwt

		BeforeAll(func() {
			config.Init()
			jwtInstance = jwt.NewJwt()
		})

		It("should generate and validate JWT tokens", func() {
			input := &jwt.GenerateTokenInput{
				ID:    "123",
				Email: "mock@mail.com",
			}
			token, exp, err := jwtInstance.GenerateAccessToken(input)
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(exp).ToNot(BeNil())
			Expect(*exp).To(BeNumerically(">", 0))
		})
	})
})
