package main

import (
	"log/slog"
	"order/internal/config"
	logCust "order/logCust"

	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	slog.Info(
		"config file param",
		slog.String("host", cfg.Db.Host),
		slog.Int("port", cfg.Db.Port),
		slog.String("User", cfg.Db.User),
		slog.String("Password", cfg.Db.Password),
		slog.String("Dbname", cfg.Db.Dbname),
	)
}

// Initializes logger
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		opts := logCust.DesignHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		}

		log = slog.New(
			logCust.NewTextDesignHandler(os.Stdout, opts),
		)
	case envDev:
		opts := logCust.DesignHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		log = slog.New(
			logCust.NewJSONDesignHandler(os.Stdout, opts),
		)
	case envProd:
		opts := logCust.DesignHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}

		log = slog.New(
			logCust.NewJSONDesignHandler(os.Stdout, opts),
		)
	}

	return log
}
