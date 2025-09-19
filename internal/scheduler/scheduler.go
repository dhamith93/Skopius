package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dhamith93/Skopius/internal/monitor"
)

// type Scheduler struct {
// 	services []monitor.Service
// 	ctx      context.Context
// 	cancel   context.CancelFunc
// 	Results  chan monitor.CheckResult
// }

type Scheduler struct {
	services map[string]monitor.Service
	cancel   map[string]context.CancelFunc
	Results  chan monitor.CheckResult
	mu       sync.Mutex
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewScheduler(services []monitor.Service) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Scheduler{
		services: make(map[string]monitor.Service),
		cancel:   make(map[string]context.CancelFunc),
		Results:  make(chan monitor.CheckResult, 100), // buffered channel
		mu:       sync.Mutex{},
		ctx:      ctx,
		cancelFn: cancel,
	}

	for _, svc := range services {
		s.services[svc.Name] = svc
	}

	return s
}

func (s *Scheduler) startService(svc monitor.Service) {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel[svc.Name] = cancel

	go func() {
		ticker := time.NewTicker(time.Duration(svc.Interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				result := svc.Check()
				s.Results <- result
			case <-ctx.Done():
				log.Printf("Stopped service: %s", svc.Name)
				return
			}
		}
	}()
}

func (s *Scheduler) Reload(newServices []monitor.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newMap := make(map[string]monitor.Service)
	for _, svc := range newServices {
		newMap[svc.Name] = svc
	}

	for name, oldSvc := range s.services {
		newSvc, exists := newMap[name]
		if !exists || oldSvc.Interval != newSvc.Interval || oldSvc.URL != newSvc.URL {
			log.Printf("Stopping service: %s", name)
			if cancel, ok := s.cancel[name]; ok {
				cancel()
				delete(s.cancel, name)
			}
			delete(s.services, name)
		}
	}

	for name, svc := range newMap {
		if _, exists := s.services[name]; !exists {
			log.Printf("Starting new service: %s", name)
			s.services[name] = svc
			s.startService(svc)
		}
	}
}

// Start launches all services with their own cancelable goroutines.
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, svc := range s.services {
		service := svc
		ctx, cancel := context.WithCancel(s.ctx)
		s.cancel[service.Name] = cancel

		go func(svc monitor.Service, ctx context.Context) {
			ticker := time.NewTicker(svc.Interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					res := svc.Check()
					s.Results <- res
				case <-ctx.Done():
					log.Printf("Stopping service: %s", svc.Name)
					return
				}
			}
		}(service, ctx)
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, cancel := range s.cancel {
		log.Printf("Cancelling service: %s", name)
		cancel()
		delete(s.cancel, name)
	}

	close(s.Results)
}
