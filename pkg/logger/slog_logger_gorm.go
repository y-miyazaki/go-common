package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

// GormSlog struct.
type GormSlog struct {
	l          *slog.Logger
	gormConfig *GormConfig
}

// GormSlogSetting sets configurations.
type GormSlogSetting struct {
	Config     *SlogConfig
	GormConfig *GormConfig
}

// NewSlogLoggerGorm creates a GORM logger backed by slog.
func NewSlogLoggerGorm(c *GormSlogSetting) *GormSlog {
	var cfg *SlogConfig
	if c == nil || c.Config == nil {
		cfg = &SlogConfig{
			Level:  LevelInfo,
			Output: os.Stdout,
			Format: "text",
		}
	} else {
		cfg = &SlogConfig{
			Level:     c.Config.Level,
			AddSource: c.Config.AddSource,
			Output:    c.Config.Output,
			Format:    c.Config.Format,
		}
	}

	slogLogger := NewSlogLogger(cfg)
	gormConfig := &GormConfig{LogLevel: Info}
	if c != nil && c.GormConfig != nil {
		gormConfig = &GormConfig{
			SlowThreshold:             c.GormConfig.SlowThreshold,
			IgnoreRecordNotFoundError: c.GormConfig.IgnoreRecordNotFoundError,
			LogLevel:                  c.GormConfig.LogLevel,
		}
	}

	return &GormSlog{
		l:          slogLogger.log,
		gormConfig: gormConfig,
	}
}

// LogMode log mode.
func (l *GormSlog) LogMode(level logger.LogLevel) logger.Interface {
	// Keep the signature stable, but avoid unused-param lint failures.
	_ = level
	return l
}

// Info print the info level log.
func (l *GormSlog) Info(ctx context.Context, msg string, data ...any) {
	_ = ctx
	if l.gormConfig.LogLevel >= Info {
		l.l.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn print the warn level log.
func (l *GormSlog) Warn(ctx context.Context, msg string, data ...any) {
	_ = ctx
	if l.gormConfig.LogLevel >= Warn {
		l.l.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error print the error level log.
func (l *GormSlog) Error(ctx context.Context, msg string, data ...any) {
	_ = ctx
	if l.gormConfig.LogLevel >= Error {
		l.l.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace print the SQL log.
func (l *GormSlog) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	_ = ctx
	level := l.gormConfig.LogLevel
	if level <= Silent {
		return
	}

	sql, rows := fc()
	duration := time.Since(begin)
	attrs := []any{
		"duration", duration.String(),
		"slowthreshold", l.gormConfig.SlowThreshold.String(),
		"sql", sql,
		"rows", rows,
	}
	if err != nil {
		attrs = append(attrs, "error", err)
	}

	switch {
	case err != nil && level >= Error && (!l.gormConfig.IgnoreRecordNotFoundError ||
		!errors.Is(err, logger.ErrRecordNotFound)):
		l.l.Error("failed SQL Query", attrs...)
	case level >= Warn && duration > l.gormConfig.SlowThreshold && l.gormConfig.SlowThreshold != 0:
		l.l.Warn(fmt.Sprintf("slow SQL Query > %v", l.gormConfig.SlowThreshold), attrs...)
	case level >= Info:
		l.l.Info("SQL Query", attrs...)
	default:
		// No logging for silent level
	}
}
