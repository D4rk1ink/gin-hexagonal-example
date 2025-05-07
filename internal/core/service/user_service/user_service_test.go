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
})
