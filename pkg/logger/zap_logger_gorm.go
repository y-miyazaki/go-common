package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// ZapInterface Gorm Logger interface.
// nolint:iface,revive,unused
type ZapInterface interface {
	LogMode(level logger.LogLevel) logger.Interface
	Info(ctx context.Context, msg string, data ...any)
	Warn(ctx context.Context, msg string, data ...any)
	Error(ctx context.Context, msg string, data ...any)
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

// LogLevelZap log level.
type LogLevelZap int

// GormZapConfig set configurations.
type GormZapConfig struct {
	// If the time specified by SlowThreshold is exceeded, it is displayed in the log as a SlowQuery.
	SlowThreshold time.Duration
	// If true is specified for IgnoreRecordNotFoundError, no log is output at the error level even if the record
	// does not exist as a search result.
	// If false, logs are output at the error level when a record does not exist as a search result.
	IgnoreRecordNotFoundError bool
	// LogLevel outputs logs above the specified level.
	LogLevel LogLevel
}

// GormZap struct.
type GormZap struct {
	l          *zap.Logger
	gormConfig *GormConfig
}

// GormSetting sets configurations.
type GormZapSetting struct {
	Config     *zap.Config
	GormConfig *GormConfig
}

// NewZapLoggerGorm func
func NewZapLoggerGorm(c *GormZapSetting) *GormZap {
	l, err := c.Config.Build()
	if err != nil {
		panic("can't create logger from zap")
	}
	return &GormZap{
		l: l,
		gormConfig: &GormConfig{
			SlowThreshold:             c.GormConfig.SlowThreshold,
			IgnoreRecordNotFoundError: c.GormConfig.IgnoreRecordNotFoundError,
			LogLevel:                  c.GormConfig.LogLevel,
		},
	}
}

// LogMode log mode(same logrus.level)
func (l *GormZap) LogMode(_ logger.LogLevel) logger.Interface { // nolint:unused
	newlogger := l
	return newlogger
}

// Info print the info level log.
func (l *GormZap) Info(_ context.Context, msg string, data ...any) { // nolint:unused
	if l.gormConfig.LogLevel >= Info {
		l.l.Sugar().Infof(msg, data...)
	}
}

// Warn print the warn level log.
func (l *GormZap) Warn(_ context.Context, msg string, data ...any) { // nolint:unused
	if l.gormConfig.LogLevel >= Warn {
		l.l.Sugar().Warnf(msg, data...)
	}
}

// Error print the error level log.
func (l *GormZap) Error(_ context.Context, msg string, data ...any) { // nolint:unused
	if l.gormConfig.LogLevel >= Error {
		l.l.Sugar().Errorf(msg, data...)
	}
}

// Trace print the SQL log.
func (l *GormZap) Trace(
	_ context.Context, // nolint:unused
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
	entry := l.l.Sugar().
		With(zap.String("duration", duration.String())).
		With(zap.String("slowthreshold", l.gormConfig.SlowThreshold.String())).
		With(zap.String("sql", sql)).
		With(zap.Int64("rows", rows))
	if err != nil {
		entry = entry.With(zap.Error(err))
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
