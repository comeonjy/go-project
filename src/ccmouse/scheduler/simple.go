package scheduler

import "ccmouse/engine"

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) Submit(request engine.Request)  {
	s.workerChan<-request
}

func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request)  {
	go func() {s.workerChan=c}()
}

