package service

import "context"

// run runs the service.
func (s *Service) run() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	defer s.cancel()
	if err := s.invokeSetupConfig(s.config); err != nil {
		return err
	}
	if err := s.invokeSetup(); err != nil {
		return err
	}
	defer s.invokeCleanup()
	return nil
}
