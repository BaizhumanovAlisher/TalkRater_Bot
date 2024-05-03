package helpers

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"log/slog"
	"os"
	"time"
)

func SetupLogger(env string, pathLogs string) *slog.Logger {
	switch env {
	case "prod":
		file, err := os.OpenFile(pathLogs, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("error in opening log file: %s", err)
		}

		return slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	case "local":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		log.Fatalf("unsupported type of logger: %s", env)
		return nil
	}
}

const opLoggerDb = "logger db"

type SlogLoggerDB struct {
	logger *slog.Logger
}

func NewSlogLoggerDB(logger *slog.Logger) *SlogLoggerDB {
	return &SlogLoggerDB{logger: logger}
}

func (s *SlogLoggerDB) LogMode(level logger.LogLevel) logger.Interface {
	return s
}

func (s *SlogLoggerDB) Info(ctx context.Context, info string, args ...interface{}) {
	fields := make([]slog.Attr, 0, len(args)+1)

	// Add the message to the fields
	fields = append(fields, slog.String("info", info))

	// Iterate over the args and add each one to the fields
	for i, arg := range args {
		fields = append(fields, slog.Any(fmt.Sprintf("args %d", i), arg))
	}

	// Log the fields
	s.logger.InfoContext(ctx, opLoggerDb, fields)
}

func (s *SlogLoggerDB) Warn(ctx context.Context, info string, args ...interface{}) {
	fields := make([]slog.Attr, 0, len(args)+1)

	// Add the message to the fields
	fields = append(fields, slog.String("info", info))

	// Iterate over the args and add each one to the fields
	for i, arg := range args {
		fields = append(fields, slog.Any(fmt.Sprintf("args %d", i), arg))
	}

	// Log the fields
	s.logger.InfoContext(ctx, opLoggerDb, fields)
}

func (s *SlogLoggerDB) Error(ctx context.Context, info string, args ...interface{}) {
	fields := make([]slog.Attr, 0, len(args)+1)

	// Add the message to the fields
	fields = append(fields, slog.String("info", info))

	// Iterate over the args and add each one to the fields
	for i, arg := range args {
		fields = append(fields, slog.Any(fmt.Sprintf("args %d", i), arg))
	}

	// Log the fields
	s.logger.InfoContext(ctx, opLoggerDb, fields)
}

func (s *SlogLoggerDB) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	duration := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		s.logger.ErrorContext(ctx, opLoggerDb,
			slog.String("error", err.Error()),
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", duration),
		)
	} else {
		s.logger.InfoContext(ctx, opLoggerDb,
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", duration),
		)
	}
}
