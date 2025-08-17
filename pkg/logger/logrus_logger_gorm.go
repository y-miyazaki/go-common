package logger

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// Interface Gorm Logger interface.
// nolint:iface,revive,unused
type Interface interface {
	LogMode(level logger.LogLevel) Interface
	Info(ctx context.Context, msg string, data ...any)
	Warn(ctx context.Context, msg string, data ...any)
	Error(ctx context.Context, msg string, data ...any)
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
	// If true is specified for IgnoreRecordNotFoundError, no log is output at the error level even if the record
	// does not exist as a search result.
	// If false, logs are output at the error level when a record does not exist as a search result.
	IgnoreRecordNotFoundError bool
	// LogLevel outputs logs above the specified level.
	LogLevel LogLevel
}

// Gorm struct.
type Gorm struct {
	l          *logrus.Logger // nolint:unused // logger instance kept for potential future use
	e          *logrus.Entry
	gormConfig *GormConfig
}

// GormSetting sets configurations.
type GormSetting struct {
	Logger     *logrus.Logger
	GormConfig *GormConfig
}

// NewLoggerGorm func
func NewLoggerGorm(c *GormSetting) *Gorm {
	var l = logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	l.SetFormatter(c.Logger.Formatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	l.SetOutput(c.Logger.Out)
	// Only log the warning severity or above.
	l.SetLevel(c.Logger.Level)
	return &Gorm{
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
func (l *Gorm) LogMode(_ logger.LogLevel) logger.Interface { // nolint:unused
	newlogger := l
	return newlogger
}

// Info print the info level log.
func (l *Gorm) Info(ctx context.Context, msg string, data ...any) {
	if l.gormConfig.LogLevel >= Info {
		l.e.WithContext(ctx).Infof(msg, data...)
	}
}

// Warn print the warn level log.
func (l *Gorm) Warn(ctx context.Context, msg string, data ...any) {
	if l.gormConfig.LogLevel >= Warn {
		l.e.WithContext(ctx).Warnf(msg, data...)
	}
}

// Error print the error level log.
func (l *Gorm) Error(ctx context.Context, msg string, data ...any) {
	if l.gormConfig.LogLevel >= Error {
		l.e.WithContext(ctx).Errorf(msg, data...)
	}
}

// Trace print the SQL log.
func (l *Gorm) Trace(
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
	case err != nil && level >= Error && (!l.gormConfig.IgnoreRecordNotFoundError ||
		!errors.Is(err, logger.ErrRecordNotFound)):
		entry.Error("failed SQL Query")
	case level >= Warn && duration > l.gormConfig.SlowThreshold && l.gormConfig.SlowThreshold != 0:
		entry.Warnf("slow SQL Query > %v", l.gormConfig.SlowThreshold)
	case level >= Info:
		entry.Info("SQL Query")
	default:
		// No logging for silent level
	}
}
