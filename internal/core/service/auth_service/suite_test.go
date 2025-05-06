package auth_service_test

import (
	"context"
	"testing"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	mock_port "github.com/D4rk1ink/gin-hexagonal-example/internal/core/port/mock"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/service"
	mock_hash "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash/mock"
	mock_jwt "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt/mock"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx         context.Context
	ctrl        *gomock.Controller
	authService port.AuthService

	mockUserRepo *mock_port.MockUserRepository
	mockJwt      *mock_jwt.MockJwt
	mockHash     *mock_hash.MockHash
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Service Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.TODO()
	ctrl = gomock.NewController(GinkgoT())

	mockUserRepo = mock_port.NewMockUserRepository(ctrl)
	mockJwt = mock_jwt.NewMockJwt(ctrl)
	mockHash = mock_hash.NewMockHash(ctrl)

	authService = service.NewAuthService(mockUserRepo, mockJwt, mockHash)
})

var _ = AfterSuite(func() {
	ctrl.Finish()
})
