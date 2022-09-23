package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
)

type IController interface {
	AddRates()
}
type ICron interface {
	Start()
	Stop() context.Context
	AddFunc(string, func()) (cron.EntryID, error)
}

// Структура планировщика
type scheduler struct {
	cr ICron
}

// Конструктор планировщика
func NewScheduler(contoller IController, cron ICron, schedule string) (*scheduler, error) {
	if schedule == "" {
		return nil, ErrNoScheduler
	}

	_, err := cron.AddFunc(schedule, contoller.AddRates)
	if err != nil {
		return nil, err
	}
	return &scheduler{cron}, nil
}

// Запуск планировщика
func (s *scheduler) Start() {
	s.cr.Start()
}

// Остановка планировщика
func (s *scheduler) Stop() {
	s.cr.Stop()
}
