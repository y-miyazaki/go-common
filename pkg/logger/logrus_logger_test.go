package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {

	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, _ = logrus.ParseLevel("Info")
	log := NewLogger(logger)
	e1 := errors.New("test1")
	e2 := errors.New("test2")
	e3 := fmt.Errorf("test3 %w", e1)
	e4 := e1

	log.WithError(e1).Error("test1")
	log.WithError(e2).Error("test2")
	log.WithError(e3).Error("test3")
	log.WithError(e4).Error("test4")

	log.Debugf("Debugf")
	log.Infof("Infof")
	log.Printf("Printf")
	log.Warnf("Warnf")
	log.Warningf("Warningf")
	log.Errorf("Errorf")
	// log.Fatalf("Fatalf")
	// log.Panicf("Panicf")
	log.Debug("Debug")
	log.Info("Info")
	log.Print("Print")
	log.Warn("Warn")
	log.Warning("Warning")
	log.Error("Error")
	// log.Fatal("Fatal")
	// log.Panic("Panic")
	log.Debugln("Debugln")
	log.Infoln("Infoln")
	log.Println("Println")
	log.Warnln("Warnln")
	log.Warningln("Warningln")
	log.Errorln("Errorln")
	// log.Fatalln("Fatalln")
	// log.Panicln("Panicln")
	log.WithField("test", "data").WithFields(logrus.Fields{"Type": "unk", "State": "oops"}).Info("test")

	e := log.GetEntry()
	assert.NotNil(t, e)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "contextKey", "contextValue")
	log.WithContext(ctx).WithContextValue("contextKey").Infof("WithContextValue")
}

func TestLogger_PanicFunctions(t *testing.T) {
	// Create a logger with a buffer to capture output
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.TextFormatter{}
	logger.Level = logrus.PanicLevel // Set to PanicLevel to ensure panic functions are called

	// Test Panicf - should panic
	assert.Panics(t, func() {
		log := NewLogger(logger)
		log.Panicf("Test panicf: %s", "message")
	})

	// Test Panic - should panic
	assert.Panics(t, func() {
		log := NewLogger(logger)
		log.Panic("Test panic")
	})

	// Test Panicln - should panic
	assert.Panics(t, func() {
		log := NewLogger(logger)
		log.Panicln("Test panicln")
	})
}

func TestLogger_FatalFunctions(t *testing.T) {
	// Note: Fatal functions call os.Exit, which terminates the program
	// In a real test environment, you might want to mock os.Exit or use a different approach
	// For now, we'll test that the functions exist and can be called (though they will exit)

	// We can't easily test Fatal functions in unit tests because they call os.Exit
	// One approach is to test that the logger is properly configured
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.TextFormatter{}
	logger.Level = logrus.FatalLevel

	log := NewLogger(logger)

	// Test that the logger has the expected level
	assert.Equal(t, logrus.FatalLevel, log.Entry.Logger.Level)

	// Fatal functions would normally exit, so we don't call them in tests
	// log.Fatalf("Test fatalf: %s", "message")
	// log.Fatal("Test fatal")
	// log.Fatalln("Test fatalln")
}
