package main

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
)

func main() {
	dep := dependency.NewDependency()

	err := http.NewHttpHandler(dep).Listen()
	if err != nil {
		panic(err)
	}
}
