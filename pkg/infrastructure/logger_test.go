package infrastructure

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, _ = logrus.ParseLevel("Info")
	log := NewLogger(logger)
	e1 := errors.WithStack(errors.New("test1"))
	e2 := errors.New("test2")
	log.WithError(e1).Errorf("test1")
	log.WithError(e2).Errorf("test2")
	log.Infof("test")

	// Infof("test")
}
