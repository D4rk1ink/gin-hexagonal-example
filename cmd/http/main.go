package main

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	err = logger.Init()
	if err != nil {
		panic(err)
	}

	dep := dependency.NewDependency()

	err = http.NewHttpHandler(dep.Service).Listen()
	if err != nil {
		panic(err)
	}
}
