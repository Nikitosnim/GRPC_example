package app

import (
	"log/slog"
	"net"
	"order/internal/grpc/crud"
	"strconv"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
}

func New(crudServise crudgrpc.Crud) *App {
	gRPCServer := crudgrpc.Register(crudServise)

	return &App{
		gRPCServer: gRPCServer,
	}
}

func (a *App) LaunchGRPCServer(port int) {
	portStr := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		slog.Error(
			"Error app.LaunchGRPCServer method net.Listen",
			slog.String("Error", err.Error()),
		)
		panic(err)
	}

	slog.Info(
		"start serwer",
		slog.String("Port", portStr),
	)

	err = a.gRPCServer.Serve(l)
	if err != nil {
		slog.Error(
			"Error app.LaunchGRPCServer",
			slog.String("Error", err.Error()),
		)
		panic(err)
	}
}

func (a *App) StopGRPCServer(port int) {
	slog.Info(
		"stopping grpc server",
		slog.String("port", ":"+strconv.Itoa(port)),
	)

	a.gRPCServer.GracefulStop()
}
