package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/neuralnorthwest/mu/worker"
)

// run runs the service.
func (s *Service) run() error {
	if s.MockMode() {
		s.logger.Info("running in mock mode")
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
	startInterruptListener(s.ctx, s.logger, s.cancel)
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

// startInterruptListener starts the interrupt listener. This registers a
// listener for SIGINT and SIGTERM signals, and starts a goroutine that
// cancels the context when a signal is received.
func startInterruptListener(ctx context.Context, logger logging.Logger, cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		logger.Info("interrupt received, shutting down")
		cancel()
	}()
}
