package integration_test

import (
	"context"
	"net/http/httptest"
	"testing"

	http_handler "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/database"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJwt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var (
	ctx      context.Context
	server   *httptest.Server
	router   *gin.Engine
	db       database.MongoDb
	jwtInfra jwt.Jwt
)

func clearDb() {
	if db != nil {
		db.GetDb().Drop(ctx)
	}
}

var _ = BeforeSuite(func() {
	ctx = context.TODO()
	config.Init()
	dep := dependency.NewDependency()
	httpHandler := http_handler.NewHttpHandler(dep)
	httpHandler.SetRouter()
	router = httpHandler.GetRouter()
	db = dep.Infrastructure.Database
	jwtInfra = dep.Infrastructure.Jwt

	db.GetDb().Drop(ctx)

	server = httptest.NewServer(router)
})

var _ = AfterSuite(func() {
	db.GetDb().Drop(ctx)
	server.Close()
})
