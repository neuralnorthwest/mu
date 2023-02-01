package service

import (
	"context"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/worker"
)

// run runs the service.
func (s *Service) run() error {
	if s.logger == nil {
		logger, err := logging.New()
		if err != nil {
			return err
		}
		s.logger = logger
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()
	if err := s.invokeConfigSetup(s.config); err != nil {
		return err
	}
	workerGroup := worker.NewGroup()
	if err := s.invokeSetup(workerGroup); err != nil {
		return err
	}
	werr := workerGroup.Run(s.ctx, s.logger)
	cerr := s.invokeCleanup()
	if werr != nil {
		return werr
	}
	if cerr != nil {
		return cerr
	}
	return nil
}
