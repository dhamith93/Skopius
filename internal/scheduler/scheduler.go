package scheduler

import (
	"time"
)

type Job func()

type Scheduler struct {
	jobs []func()
}

func New() *Scheduler {
	return &Scheduler{jobs: []func(){}}
}

func (s *Scheduler) Every(interval int, job Job) {
	wrapped := func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		go func() {
			for range ticker.C {
				job()
			}
		}()
	}
	s.jobs = append(s.jobs, wrapped)
}

func (s *Scheduler) Start() {
	for _, j := range s.jobs {
		j()
	}
	select {} // block forever
}
