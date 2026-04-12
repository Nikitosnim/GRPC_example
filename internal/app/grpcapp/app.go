package grpcapp

import (
	"log/slog"
	"net"
	"order/internal/grpc/crudgrpc"
	"strconv"

	"google.golang.org/grpc"
)

type GrpcApp struct {
	gRPCServer *grpc.Server
}

func New(crudService crudgrpc.Crud) *GrpcApp {
	gRPCServer := crudgrpc.Register(crudService)

	return &GrpcApp{
		gRPCServer: gRPCServer,
	}
}

func (a *GrpcApp) StartGRPCServer(port int) {
	portStr := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		slog.Error(
			"Error app.StartGRPCServer method net.Listen",
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
			"Error app.StartGRPCServer",
			slog.String("Error", err.Error()),
		)
		panic(err)
	}
}

func (a *GrpcApp) StopGRPCServer(port int) {
	slog.Info(
		"stopping grpc server",
		slog.String("Port", ":"+strconv.Itoa(port)),
	)

	a.gRPCServer.GracefulStop()
}
