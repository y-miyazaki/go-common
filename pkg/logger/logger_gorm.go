package logger

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// Interface Gorm Logger interface.
type Interface interface {
	LogMode(logger.LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

// LogLevel log level.
type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

// GormConfig set configurations.
type GormConfig struct {
	// If the time specified by SlowThreshold is exceeded, it is displayed in the log as a SlowQuery.
	SlowThreshold time.Duration
	// If true is specified for IgnoreRecordNotFoundError, no log is output at the error level even if the record does not exist as a search result.
	// If false, logs are output at the error level when a record does not exist as a search result.
	IgnoreRecordNotFoundError bool
	// LogLevel outputs logs above the specified level.
	LogLevel LogLevel
}

// LoggerGorm struct.
type LoggerGorm struct {
	l          *logrus.Logger
	e          *logrus.Entry
	gormConfig *GormConfig
}

// LoggerGormConfig sets configurations.
type LoggerGormConfig struct {
	Logger     *logrus.Logger
	GormConfig *GormConfig
}

// NewLoggerGorm func
func NewLoggerGorm(c *LoggerGormConfig) *LoggerGorm {
	var l = logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	l.SetFormatter(c.Logger.Formatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	l.SetOutput(c.Logger.Out)
	// Only log the warning severity or above.
	l.SetLevel(c.Logger.Level)
	return &LoggerGorm{
		l: l,
		e: l.WithFields(logrus.Fields{}),
		gormConfig: &GormConfig{
			SlowThreshold:             c.GormConfig.SlowThreshold,
			IgnoreRecordNotFoundError: c.GormConfig.IgnoreRecordNotFoundError,
			LogLevel:                  c.GormConfig.LogLevel,
		},
	}
}

// LogMode log mode(same logrus.level)
func (l *LoggerGorm) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := l
	return newlogger
}

// Info print the info level log.
func (l *LoggerGorm) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.gormConfig.LogLevel >= Info {
		l.e.WithContext(ctx).Infof(msg, data...)
	}
}

// Warn print the warn level log.
func (l *LoggerGorm) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.gormConfig.LogLevel >= Warn {
		l.e.WithContext(ctx).Warnf(msg, data...)
	}
}

// Error print the error level log.
func (l *LoggerGorm) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.gormConfig.LogLevel >= Error {
		l.e.WithContext(ctx).Errorf(msg, data...)
	}
}

// Trace print the SQL log.
func (l *LoggerGorm) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	level := l.gormConfig.LogLevel
	if level <= Silent {
		return
	}
	sql, rows := fc()
	duration := time.Since(begin)
	entry := l.e.
		WithContext(ctx).
		WithField("duration", duration.String()).
		WithField("sql", sql).
		WithField("rows", rows)
	if err != nil {
		entry = entry.WithField("error", err)
	}

	switch {
	case err != nil && level >= Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.gormConfig.IgnoreRecordNotFoundError):
		entry.Error("failed SQL Query")
	case level >= Warn && duration > l.gormConfig.SlowThreshold && l.gormConfig.SlowThreshold != 0:
		entry.Warnf("slow SQL Query > %v", l.gormConfig.SlowThreshold)
	case level >= Info:
		entry.Info("SQL Query")
	}
}
