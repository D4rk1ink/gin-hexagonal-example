package user_service_test

import (
	"errors"
	"time"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	time_util "github.com/D4rk1ink/gin-hexagonal-example/internal/util/time"
	"github.com/guregu/null"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Service", Label("Service"), func() {
	var users []*domain.User
	var mockUpdatedTime time.Time

	BeforeEach(func() {
		mockTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
		mockUpdatedTime, _ = time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")
		time_util.Now = func() time.Time {
			return mockUpdatedTime
		}
		users = []*domain.User{
			{
				ID:        "1",
				Name:      "John Doe",
				Email:     "mock.1@email.com",
				Password:  "password",
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
			},
			{
				ID:        "2",
				Name:      "Kim Doe",
				Email:     "mock.2@email.com",
				Password:  "password",
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
			},
		}
	})

	Context("Count", func() {
		It("should return user count", func() {
			mockUserRepo.
				EXPECT().
				Count(ctx).
				Return(int64(len(users)), nil)
			result, err := userService.Count(ctx)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result).To(Equal(int64(len(users))))
		})
		It("should return error when user repo Count return error", func() {
			mockUserRepo.
				EXPECT().
				Count(ctx).
				Return(int64(0), errors.New("error"))
			result, err := userService.Count(ctx)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeZero())
			Expect(err.Error()).To(Equal(errors.New("error").Error()))
		})
	})

	Context("GetAll", func() {
		It("should return all users", func() {
			mockUserRepo.
				EXPECT().
				GetAll(ctx).
				Return(users, nil)

			result, err := userService.GetAll(ctx)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result).To(HaveLen(2))
		})
		It("should return empty if no any user", func() {
			mockUserRepo.
				EXPECT().
				GetAll(ctx).
				Return([]*domain.User{}, nil)

			result, err := userService.GetAll(ctx)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result).To(BeEmpty())
		})
	})

	Context("GetById", func() {
		It("should return user by id", func() {
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(users[0], nil)

			result, err := userService.GetById(ctx, users[0].ID)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result.ID).To(Equal(users[0].ID))
		})
		It("should return nil if user not found", func() {
			mockUserRepo.
				EXPECT().
				GetById(ctx, "invalid_id").
				Return(nil, nil)

			result, err := userService.GetById(ctx, "invalid_id")

			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})

	Context("Create", func() {
		It("should return user domain", func() {
			user := users[0]
			payload := dto.UserCreateDto{
				Name:            user.Name,
				Email:           user.Email,
				Password:        user.Password,
				ConfirmPassword: user.Password,
			}
			expectedUserDomain, _ := domain.NewUser(user.Name, user.Email, "hashed_password")

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(nil, nil)
			mockHash.
				EXPECT().
				HashPassword(ctx, payload.Password).
				Return(null.StringFrom("hashed_password").Ptr(), nil)
			mockUserRepo.
				EXPECT().
				Create(ctx, expectedUserDomain).
				Return(&user.ID, nil)
			mockUserRepo.
				EXPECT().
				GetById(ctx, user.ID).
				Return(user, nil)

			result, err := userService.Create(ctx, payload)

			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
			Expect(result).To(Equal(user))
		})
		It("should return error when email already exists", func() {
			user := users[0]
			payload := dto.UserCreateDto{
				Name:            user.Name,
				Email:           user.Email,
				Password:        user.Password,
				ConfirmPassword: user.Password,
			}

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(user, nil)

			result, err := userService.Create(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomErrorInterface).GetCode()).To(Equal(custom_error.ErrAuthEmailAlreadyExists))
		})
		It("should return error when password and confirm password not match", func() {
			user := users[0]
			payload := dto.UserCreateDto{
				Name:            user.Name,
				Email:           user.Email,
				Password:        "password",
				ConfirmPassword: "invalid_password",
			}

			result, err := userService.Create(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomErrorInterface).GetCode()).To(Equal(custom_error.ErrAuthInvalidConfirmPassword))
		})
		It("should return error when hash password failed", func() {
			user := users[0]
			payload := dto.UserCreateDto{
				Name:            user.Name,
				Email:           user.Email,
				Password:        user.Password,
				ConfirmPassword: user.Password,
			}

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(nil, nil)
			mockHash.
				EXPECT().
				HashPassword(ctx, payload.Password).
				Return(nil, errors.New("hash error"))

			result, err := userService.Create(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.Error()).To(Equal(errors.New("hash error").Error()))
		})
		It("should return error when create user failed", func() {
			user := users[0]
			payload := dto.UserCreateDto{
				Name:            user.Name,
				Email:           user.Email,
				Password:        user.Password,
				ConfirmPassword: user.Password,
			}
			expectedUserDomain, _ := domain.NewUser(user.Name, user.Email, "hashed_password")

			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, payload.Email).
				Return(nil, nil)
			mockHash.
				EXPECT().
				HashPassword(ctx, payload.Password).
				Return(null.StringFrom("hashed_password").Ptr(), nil)
			mockUserRepo.
				EXPECT().
				Create(ctx, expectedUserDomain).
				Return(nil, errors.New("create error"))

			result, err := userService.Create(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.Error()).To(Equal(errors.New("create error").Error()))
		})
	})

	Context("Update", func() {
		It("should return updated user", func() {
			payload := dto.UserUpdateDto{
				ID:    users[0].ID,
				Name:  null.StringFrom("John Doe Updated").Ptr(),
				Email: null.StringFrom("mock.update@email.com").Ptr(),
			}
			updatedUser := *users[0]
			updatedUser.SetName(*payload.Name)
			updatedUser.SetEmail(*payload.Email)
			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, *payload.Email).
				Return(nil, nil)
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(users[0], nil)
			mockUserRepo.
				EXPECT().
				Update(ctx, users[0]).
				Return(nil)

			result, err := userService.Update(ctx, payload)

			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
			Expect(result.ID).To(Equal(users[0].ID))
			Expect(result.Name).To(Equal(*payload.Name))
			Expect(result.Email).To(Equal(*payload.Email))
			Expect(result.UpdatedAt).To(Equal(mockUpdatedTime))
		})
		It("should return error if input is invalid", func() {
			payload := dto.UserUpdateDto{
				ID:    users[0].ID,
				Name:  nil,
				Email: nil,
			}

			result, err := userService.Update(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomError).GetCode()).To(Equal(custom_error.ErrUserInvalidateUpdateInput))
		})
		It("should return error if email already exists", func() {
			payload := dto.UserUpdateDto{
				ID:    users[0].ID,
				Name:  null.StringFrom(users[0].Name).Ptr(),
				Email: null.StringFrom(users[0].Email).Ptr(),
			}
			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, *payload.Email).
				Return(users[0], nil)

			result, err := userService.Update(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomError).GetCode()).To(Equal(custom_error.ErrUserEmailAlreadyExists))
		})
		It("should return error if user not found", func() {
			payload := dto.UserUpdateDto{
				ID:    users[0].ID,
				Name:  null.StringFrom(users[0].Name).Ptr(),
				Email: null.StringFrom(users[0].Email).Ptr(),
			}
			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, *payload.Email).
				Return(nil, nil)
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(nil, nil)

			result, err := userService.Update(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
			Expect(err.(custom_error.CustomError).GetCode()).To(Equal(custom_error.ErrUserNotFound))
		})
		It("should return error if user repo GetById return error", func() {
			payload := dto.UserUpdateDto{
				ID:    users[0].ID,
				Name:  null.StringFrom(users[0].Name).Ptr(),
				Email: null.StringFrom(users[0].Email).Ptr(),
			}
			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, *payload.Email).
				Return(nil, nil)
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(nil, errors.New("error"))

			result, err := userService.Update(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
		It("should return error if user repo Update return error", func() {
			payload := dto.UserUpdateDto{
				ID:    users[0].ID,
				Name:  null.StringFrom("new name").Ptr(),
				Email: null.StringFrom("new@email.com").Ptr(),
			}
			mockUserRepo.
				EXPECT().
				GetByEmail(ctx, *payload.Email).
				Return(nil, nil)
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(users[0], nil)
			mockUserRepo.
				EXPECT().
				Update(ctx, users[0]).
				Return(errors.New("error"))

			result, err := userService.Update(ctx, payload)

			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})

	Context("Delete", func() {
		It("should return nil if user deleted successfully", func() {
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(users[0], nil)
			mockUserRepo.
				EXPECT().
				Delete(ctx, users[0].ID).
				Return(nil)

			err := userService.Delete(ctx, users[0].ID)

			Expect(err).To(BeNil())
		})
		It("should return error if user not found", func() {
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(nil, nil)

			err := userService.Delete(ctx, users[0].ID)

			Expect(err).To(HaveOccurred())
		})
		It("should return error if cannot get user by id", func() {
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(nil, errors.New("error"))

			err := userService.Delete(ctx, users[0].ID)

			Expect(err).To(HaveOccurred())
		})
		It("should return error if user repo Delete return error", func() {
			mockUserRepo.
				EXPECT().
				GetById(ctx, users[0].ID).
				Return(users[0], nil)
			mockUserRepo.
				EXPECT().
				Delete(ctx, users[0].ID).
				Return(errors.New("error"))

			err := userService.Delete(ctx, users[0].ID)

			Expect(err).To(HaveOccurred())
		})
	})
})
