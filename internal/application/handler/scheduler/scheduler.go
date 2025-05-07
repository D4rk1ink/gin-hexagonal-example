package scheduler_handler

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	LogCountUsers(ctx context.Context) (gocron.Job, error)
}

type scheduler struct {
	userService port.UserService
	cron        gocron.Scheduler
}

func NewScheduler(dep *dependency.Dependency) Scheduler {
	cron, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	return &scheduler{
		userService: dep.Service.UserService,
		cron:        cron,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	_, err := s.LogCountUsers(ctx)
	if err != nil {
		return err
	}

	s.cron.Start()

	return nil
}

func (s *scheduler) Shutdown(ctx context.Context) error {
	s.cron.Shutdown()
	return nil
}
