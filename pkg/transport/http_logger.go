package transport

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/logger"
)

// TransportHTTPLoggerType defines the transport type.
type TransportHTTPLoggerType string

const (
	// TransportHTTPLoggerTypeExternal defines the external type.
	TransportHTTPLoggerTypeExternal TransportHTTPLoggerType = "external"
	// TransportHTTPLoggerTypeInternal defines the internal type.
	TransportHTTPLoggerTypeInternal TransportHTTPLoggerType = "internal"
)

// TransportHTTPLogger struct.
type TransportHTTPLogger struct {
	http.RoundTripper
	logger *logger.Logger
	Type   TransportHTTPLoggerType
}

// NewTransportHTTPLogger get http.RoundTripper.
func NewTransportHTTPLogger(
	logger *logger.Logger,
	transportType TransportHTTPLoggerType,
) http.RoundTripper {
	return &TransportHTTPLogger{
		http.DefaultTransport,
		logger,
		transportType,
	}
}

// RoundTrip logs transparently.
func (t TransportHTTPLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	timeBefore := time.Now()
	response, err := t.RoundTripper.RoundTrip(req)
	timeAfter := time.Now()

	logger := t.logger.WithFields(logrus.Fields{
		"url":           req.URL.String(),
		"method":        req.Method,
		"protocol":      req.Proto,
		"duration":      timeAfter.Sub(timeBefore).String(),
		"transportType": t.Type,
	})
	if response != nil {
		logger = logger.WithField("status", response.StatusCode)
	}
	if err != nil || response.StatusCode/100 >= 4 {
		logger.WithError(err).Error()
	} else {
		logger.Info()
	}
	return response, err
}
