package auth_service_test

import (
	"errors"
	"time"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"github.com/guregu/null"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Service", Label("Service"), func() {
	var user *domain.User

	BeforeEach(func() {
		mockTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
		user = &domain.User{
			ID:        "1",
			Name:      "John Doe",
			Email:     "mock@email.com",
			Password:  "password",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		}
	})

	Context("Login", func() {
		It("should return access token", func() {
			payload := dto.CredentialDto{
				Email:    user.Email,
				Password: user.Password,
			}

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(user, nil)
			mockHash.
				EXPECT().
				ComparePassword(ctx, user.Password, user.Password).
				Return(nil)
			mockJwt.
				EXPECT().
				GenerateAccessToken(&jwt.GenerateTokenInput{
					ID:    user.ID,
					Email: user.Email,
				}).
				Return(null.StringFrom("jwt_access_token").Ptr(), null.IntFrom(3600).Ptr(), nil)

			result, err := authService.Login(ctx, payload)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.AccessToken).ToNot(BeEmpty())
			Expect(result.TokenType).To(Equal("Bearer"))
			Expect(result.ExpiresIn).To(BeNumerically(">", 0))
		})
		It("should return error when user not found", func() {
			payload := dto.CredentialDto{
				Email:    user.Email,
				Password: user.Password,
			}

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(nil, nil)

			result, err := authService.Login(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomErrorInterface).GetCode()).To(Equal(custom_error.ErrAuthInvalidCredentials))
		})
		It("should return error when password is invalid", func() {
			payload := dto.CredentialDto{
				Email:    user.Email,
				Password: "invalid_password",
			}

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(user, nil)
			mockHash.
				EXPECT().
				ComparePassword(ctx, payload.Password, user.Password).
				Return(errors.New("invalid password"))

			result, err := authService.Login(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomErrorInterface).GetCode()).To(Equal(custom_error.ErrAuthInvalidCredentials))
		})
		It("should return error when jwt generate access token failed", func() {
			payload := dto.CredentialDto{
				Email:    user.Email,
				Password: user.Password,
			}

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(user, nil)
			mockHash.
				EXPECT().
				ComparePassword(ctx, payload.Password, user.Password).
				Return(nil)
			mockJwt.
				EXPECT().
				GenerateAccessToken(&jwt.GenerateTokenInput{
					ID:    user.ID,
					Email: user.Email,
				}).
				Return(nil, nil, errors.New("jwt error"))

			result, err := authService.Login(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.Error()).To(Equal(errors.New("jwt error").Error()))
		})
	})
})
