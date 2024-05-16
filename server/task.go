package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/aksbuzz/library-project/store"
	"github.com/go-co-op/gocron/v2"
)

type TaskScheduler struct {
	logger    *slog.Logger
	store     *store.Store
	scheduler gocron.Scheduler
}

func NewTaskScheduler(store *store.Store, logger *slog.Logger) *TaskScheduler {
	return &TaskScheduler{store: store, logger: logger.With("service", "task-scheduler")}
}

func (t *TaskScheduler) Start(ctx context.Context) error {
	t.logger.Info("Starting task scheduler")
	options := []gocron.SchedulerOption{
		gocron.WithLocation(time.UTC),
	}

	// create a scheduler
	s, err := gocron.NewScheduler(options...)
	if err != nil {
		return err
	}
	t.scheduler = s

	// add jobs to scheduler

	// calculate overdue fees daily at 1am
	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(1, 0, 0))),
		gocron.NewTask(
			func() {
				_, err := t.store.DB.Exec(ctx, "CALL calculate_overdue_fees()")
				if err != nil {
					t.logger.Error("error calculating overdue fees", slog.Any("error", err))
				}
			},
		),
	)
	if err != nil {
		return err
	}

	// send overdue fees email weekly
	_, err = s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(gocron.NewAtTime(13, 0, 0))),
		gocron.NewTask(
			func() {
				t.logger.Info("Sending overdue fees email")
			},
		),
	)
	if err != nil {
		return err
	}

	t.scheduler.Start()

	select {
	case <-time.After(time.Minute):
	}

	return nil
}

func (t *TaskScheduler) Stop(ctx context.Context) error {
	err := t.scheduler.Shutdown()
	if err != nil {
		return err
	}

	t.logger.Info("Task scheduler stopped")

	return nil
}
