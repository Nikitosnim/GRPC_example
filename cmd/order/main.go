package main

import (
	"context"
	"log/slog"
	"order/internal/app"
	"order/internal/config"
	logCustom "order/logCust"
	"os/signal"
	"syscall"

	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	slog.Info(
		"config file param",
		slog.String("host", cfg.Db.Host),
		slog.Int("port", cfg.Db.Port),
		slog.String("User", cfg.Db.User),
		slog.String("Dbname", cfg.Db.Dbname),
	)

	appliction := app.New(ctx, cfg)
	go appliction.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	slog.Info(
		"stopping application",
		slog.String("Signal:", sign.String()),
	)

	appliction.Stop()

	slog.Info("application stop")

}

// Initializes logger
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		opts := logCustom.DesignHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		}

		log = slog.New(
			logCustom.NewTextDesignHandler(os.Stdout, opts),
		)
	case envDev:
		opts := logCustom.DesignHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		log = slog.New(
			logCustom.NewJSONDesignHandler(os.Stdout, opts),
		)
	case envProd:
		opts := logCustom.DesignHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}

		log = slog.New(
			logCustom.NewJSONDesignHandler(os.Stdout, opts),
		)
	default:
		panic("Environment is missing to logger main.setupLogger")
	}

	return log
}
