package jwt_test

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT", Label("Infrastructure"), func() {
	var generateTokenOptions *jwt.GenerateTokenOptions
	var validateTokenOptions *jwt.ValidateTokenOptions

	BeforeEach(func() {
		generateTokenOptions = &jwt.GenerateTokenOptions{
			Secret:   config.Config.Jwt.Secret,
			Duration: config.Config.Jwt.ExpiresIn,
		}
		validateTokenOptions = &jwt.ValidateTokenOptions{
			Secret: config.Config.Jwt.Secret,
		}
	})

	Context("GenerateTokenWithOptions", func() {
		It("should generate and validate JWT tokens", func() {
			input := &jwt.GenerateTokenInput{
				ID:    "123",
				Email: "mock@mail.com",
			}
			token, exp, err := jwtInstance.GenerateTokenWithOptions(input, generateTokenOptions)

			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeNil())
			Expect(exp).ToNot(BeNil())
			Expect(*exp).To(BeNumerically(">", 0))
		})
	})

	Context("ValidateTokenWithOptions", func() {
		It("should validate and return claim correctly", func() {
			input := &jwt.GenerateTokenInput{
				ID:    "123",
				Email: "mock@mail.com",
			}
			token, _, _ := jwtInstance.GenerateTokenWithOptions(input, generateTokenOptions)
			claims, err := jwtInstance.ValidateTokenWithOptions(*token, validateTokenOptions)

			Expect(err).ToNot(HaveOccurred())
			Expect(claims).ToNot(BeNil())
			Expect(claims).To(BeAssignableToTypeOf(&jwt.TokenPayload{}))
			Expect(claims.ID).To(Equal(input.ID))
			Expect(claims.Email).To(Equal(input.Email))
		})

		It("should return error when token is invalid format", func() {
			input := &jwt.GenerateTokenInput{
				ID:    "123",
				Email: "mock@mail.com",
			}
			token, _, _ := jwtInstance.GenerateTokenWithOptions(input, generateTokenOptions)
			invalidToken := *token + "invalid"
			claims, err := jwtInstance.ValidateTokenWithOptions(invalidToken, validateTokenOptions)

			Expect(err).To(HaveOccurred())
			Expect(claims).To(BeNil())
		})

		It("should return error when token is invalid signature", func() {
			input := &jwt.GenerateTokenInput{
				ID:    "123",
				Email: "mock@mail.com",
			}
			mockGenerateTokenOptions := &jwt.GenerateTokenOptions{
				Secret:   "invalidsecret",
				Duration: config.Config.Jwt.ExpiresIn,
			}
			token, _, _ := jwtInstance.GenerateTokenWithOptions(input, mockGenerateTokenOptions)
			claims, err := jwtInstance.ValidateTokenWithOptions(*token, validateTokenOptions)

			Expect(err).To(HaveOccurred())
			Expect(claims).To(BeNil())
		})
	})
})
