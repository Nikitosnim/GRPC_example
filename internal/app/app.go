package app

import (
	"context"
	"fmt"
	"order/internal/app/grpcapp"
	"order/internal/config"
	"order/internal/services/order"
	"order/internal/storage/postgresql"
)

type App struct {
	ctx     context.Context
	cfg     *config.Config
	grpcApp *grpcapp.GrpcApp
	storage order.Storage
}

func New(ctx context.Context, cfg *config.Config) *App {
	return &App{
		ctx: ctx,
		cfg: cfg,
	}
}

func (a *App) MustRun() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		a.cfg.Db.User, a.cfg.Db.Password,
		a.cfg.Db.Host,
		a.cfg.Db.Port,
		a.cfg.Db.Dbname,
	)

	// storage layer
	storage, err := postgresql.New(a.ctx, connStr)
	if err != nil {
		panic("Error when connecting db" + err.Error())
	}
	a.storage = storage

	// servise layer
	crudService := order.New(storage)

	grpcApp := grpcapp.New(crudService)
	a.grpcApp = grpcApp

	grpcApp.StartGRPCServer(a.cfg.GRPC.Port)
}

func (a *App) Stop() {
	if a.grpcApp != nil {
		a.grpcApp.StopGRPCServer(a.cfg.GRPC.Port)
	}

	if a.storage != nil {
		a.storage.Stop()
	}
}
