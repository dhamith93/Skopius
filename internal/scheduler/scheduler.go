package scheduler

import (
	"context"
	"time"

	"github.com/dhamith93/Skopius/internal/monitor"
)

type Scheduler struct {
	services []monitor.Service
	ctx      context.Context
	cancel   context.CancelFunc
	Results  chan monitor.CheckResult
}

func NewScheduler(services []monitor.Service) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		services: services,
		ctx:      ctx,
		cancel:   cancel,
		Results:  make(chan monitor.CheckResult, 100), // buffered channel
	}
}

func (s *Scheduler) Start() {
	for _, svc := range s.services {
		service := svc
		go func() {
			ticker := time.NewTicker(service.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					res := service.Check()
					s.Results <- res
				case <-s.ctx.Done():
					return
				}
			}
		}()
	}
}

func (s *Scheduler) Stop() {
	s.cancel()
	close(s.Results)
}
