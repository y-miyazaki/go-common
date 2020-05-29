package infrastructure

import (
	"github.com/sirupsen/logrus"
)

// NewLogger returns an instance of logger
func NewLogger(logger *logrus.Logger) (*logrus.Logger, error) {
	var log = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(logger.Formatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(logger.Out)
	// Only log the warning severity or above.
	log.SetLevel(logger.Level)
	return log, nil
}
