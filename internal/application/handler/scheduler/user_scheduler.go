package scheduler_handler

import (
	"context"
	"time"

	"fmt"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/go-co-op/gocron/v2"
)

func (s *scheduler) LogCountUsers(ctx context.Context) (gocron.Job, error) {
	jobDuration := gocron.DurationJob(time.Duration(10) * time.Second)
	task := gocron.NewTask(func() {
		count, err := s.userService.Count(ctx)
		if err != nil {
			logger.Error("Error counting users: " + err.Error())
			return
		}

		logger.Info("Count of users: " + fmt.Sprint(count))
	})

	return s.cron.NewJob(jobDuration, task)
}
