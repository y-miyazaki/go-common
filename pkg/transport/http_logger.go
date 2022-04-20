package transport

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// HTTPLoggerType defines the transport type.
type HTTPLoggerType string

const (
	// HTTPLoggerTypeExternal defines the external type.
	HTTPLoggerTypeExternal HTTPLoggerType = "external"
	// HTTPLoggerTypeInternal defines the internal type.
	HTTPLoggerTypeInternal HTTPLoggerType = "internal"
)

// HTTPLogger struct.
type HTTPLogger struct {
	http.RoundTripper
	logger *logger.Logger
	Type   HTTPLoggerType
}

// NewTransportHTTPLogger get http.RoundTripper.
func NewTransportHTTPLogger(
	l *logger.Logger,
	transportType HTTPLoggerType,
) http.RoundTripper {
	return &HTTPLogger{
		http.DefaultTransport,
		l,
		transportType,
	}
}

// RoundTrip logs transparently.
func (t HTTPLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	timeBefore := time.Now()
	response, err := t.RoundTripper.RoundTrip(req)
	timeAfter := time.Now()

	log := t.logger.WithFields(logrus.Fields{
		"url":           req.URL.String(),
		"method":        req.Method,
		"protocol":      req.Proto,
		"duration":      timeAfter.Sub(timeBefore).String(),
		"transportType": t.Type,
	})
	if response != nil {
		log = log.WithField("status", response.StatusCode)
	}
	if err != nil || response.StatusCode/100 >= 4 {
		log.WithError(err).Error()
	} else {
		log.Info()
	}
	return response, err
}
