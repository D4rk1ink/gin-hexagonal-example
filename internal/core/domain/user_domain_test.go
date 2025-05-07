package domain_test

import (
	"time"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	time_util "github.com/D4rk1ink/gin-hexagonal-example/internal/util/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Domain", Label("Domain"), func() {
	var mockTime time.Time

	BeforeEach(func() {
		mockTime, _ = time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
		time_util.Now = func() time.Time {
			return mockTime
		}
	})

	Context("User", func() {
		It("should create a new user", func() {
			user, err := domain.NewUser("John Doe", "john@email.com", "password")

			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())
			Expect(user.Name).To(Equal("John Doe"))
			Expect(user.Email).To(Equal("john@email.com"))
			Expect(user.Password).To(Equal("password"))
			Expect(user.CreatedAt).To(Equal(mockTime))
			Expect(user.UpdatedAt).To(Equal(mockTime))
		})
		It("should return error for invalid email format", func() {
			user, err := domain.NewUser("John Doe", "invalid-email", "password")

			Expect(err).NotTo(BeNil())
			Expect(user).To(BeNil())
			Expect(err.(custom_error.CustomError).Code).To(Equal(custom_error.ErrUserInvalidEmailFormat))
		})
		It("should update an existing user", func() {
			_mockTime, _ := time.Parse(time.RFC3339, "2023-01-02T00:00:00Z")
			time_util.Now = func() time.Time {
				return _mockTime
			}
			user, err := domain.NewUser("John Doe", "john@email.com", "password")
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())

			user.SetName("Jane Doe New")
			err = user.SetEmail("new_john@email.com")
			Expect(err).To(BeNil())
			Expect(user.Name).To(Equal("Jane Doe New"))
			Expect(user.Email).To(Equal("new_john@email.com"))
			Expect(user.UpdatedAt).To(Equal(_mockTime))
		})
		It("should return error for invalid email format on update", func() {
			user, err := domain.NewUser("John Doe", "john@email.com", "password")
			Expect(err).To(BeNil())
			Expect(user).NotTo(BeNil())

			err = user.SetEmail("new_john_email.com")
			Expect(err).To(HaveOccurred())
			Expect(err.(custom_error.CustomError).Code).To(Equal(custom_error.ErrUserInvalidEmailFormat))
		})
	})
})
