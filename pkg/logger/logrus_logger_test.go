package logger

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {

	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, _ = logrus.ParseLevel("Info")
	log := NewLogger(logger)
	e1 := errors.WithStack(errors.New("test1"))
	e2 := errors.New("test2")
	e3 := fmt.Errorf("test3 %s", e1)
	e4 := errors.WithStack(e1)

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
