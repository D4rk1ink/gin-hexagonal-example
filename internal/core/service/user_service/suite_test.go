package user_service_test

import (
	"context"
	"testing"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	mock_port "github.com/D4rk1ink/gin-hexagonal-example/internal/core/port/mock"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/service/user_service"
	mock_hash "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx         context.Context
	ctrl        *gomock.Controller
	userService port.UserService

	mockUserRepo *mock_port.MockUserRepository
	mockHash     *mock_hash.MockHash
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Service Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.TODO()
	ctrl = gomock.NewController(GinkgoT())

	mockUserRepo = mock_port.NewMockUserRepository(ctrl)
	mockHash = mock_hash.NewMockHash(ctrl)

	userService = user_service.NewUserService(mockUserRepo, mockHash)
})

var _ = AfterSuite(func() {
	ctrl.Finish()
})
