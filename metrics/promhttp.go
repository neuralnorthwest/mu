package metrics

import (
	ht "net/http"

	"github.com/neuralnorthwest/mu/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns an ht.Handler that serves the metrics. The provided logger
// will be used to log any errors that occur while serving the metrics.
func (m *metrics) Handler(logger logging.Logger) ht.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{
		ErrorLog:      logging.NewAdapter(logger, logging.AdaptedLevelError),
		ErrorHandling: promhttp.ContinueOnError,
	})
}
