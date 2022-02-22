package transport

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
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
	Entry *logrus.Entry
	Type  TransportHTTPLoggerType
}

// NewTransportHTTPLogger get http.RoundTripper.
func NewTransportHTTPLogger(
	entry *logrus.Entry,
	transportType TransportHTTPLoggerType,
) http.RoundTripper {
	return &TransportHTTPLogger{
		http.DefaultTransport,
		entry,
		transportType,
	}
}

// RoundTrip logs transparently.
func (t TransportHTTPLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	timeBefore := time.Now()
	response, err := t.RoundTripper.RoundTrip(req)
	timeAfter := time.Now()

	e := t.Entry.WithFields(logrus.Fields{
		"url":           req.URL.String(),
		"method":        req.Method,
		"protocol":      req.Proto,
		"duration":      timeAfter.Sub(timeBefore).String(),
		"transportType": t.Type,
		"status":        response.Status,
	})

	if err != nil || response.StatusCode/100 >= 4 {
		e.WithError(err).Error()
	} else {
		e.Info()
	}
	return response, err
}
