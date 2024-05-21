package helpers

import (
	"context"
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
		file, err := os.OpenFile(pathLogs+string(os.PathSeparator)+"logs.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("error in opening log file: %s", err)
		}

		return slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), &slog.HandlerOptions{
			Level: slog.LevelInfo,
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

const op = "logger.db"

type SlogLoggerDB struct {
	logger *slog.Logger
}

func NewSlogLoggerDB(logger *slog.Logger) *SlogLoggerDB {
	return &SlogLoggerDB{logger: logger}
}

func (s *SlogLoggerDB) LogMode(_ logger.LogLevel) logger.Interface {
	return s
}

func (s *SlogLoggerDB) Info(ctx context.Context, info string, args ...interface{}) {
	s.logger.InfoContext(ctx, info,
		slog.String("op", op),
		slog.Any("args", args),
	)
}

func (s *SlogLoggerDB) Warn(ctx context.Context, info string, args ...interface{}) {
	s.logger.WarnContext(ctx, info,
		slog.String("op", op),
		slog.Any("args", args),
	)
}

func (s *SlogLoggerDB) Error(ctx context.Context, info string, args ...interface{}) {
	s.logger.ErrorContext(ctx, info,
		slog.String("op", op),
		slog.Any("args", args),
	)
}

func (s *SlogLoggerDB) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	duration := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(),
			slog.String("op", op),
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", duration),
		)
	} else {
		s.logger.InfoContext(ctx, "",
			slog.String("op", op),
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", duration),
		)
	}
}
