package helpers

import (
	"io"
	"log"
	"log/slog"
	"os"
)

const ()

func SetupLogger(env string) *slog.Logger {
	switch env {
	case "prod":
		pathFile := os.Getenv("PATH_LOGS")
		file, err := os.OpenFile(pathFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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
