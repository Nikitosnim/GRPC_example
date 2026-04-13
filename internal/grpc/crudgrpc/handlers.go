package crudgrpc

import (
	"context"
	"log/slog"
	"order/internal/domain/models"

	orderv1 "github.com/Nikitosnim/protos/gen/go/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	orderv1.UnimplementedCRUDServer
	crud Crud
}

// Register GRPC Server
func Register(crud Crud) *grpc.Server {
	s := grpc.NewServer()
	srv := &serverAPI{crud: crud}
	orderv1.RegisterCRUDServer(s, srv)

	return s
}

// Inetrfase CRUD servise
//
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=handlers.go -package=mocks -destination=mocks/crud_mock.go
type Crud interface {
	Create(
		ctx context.Context,
		item string,
		quan int32,
	) (string, error)
	Get(
		ctx context.Context,
		id string,
	) (*models.Order, error)
	Update(
		ctx context.Context,
		id string,
		item string,
		quan int32,
	) (*models.Order, error)
	Delete(
		ctx context.Context,
		id string,
	) (bool, error)
	List(ctx context.Context) ([]*models.Order, error)
}

// Hanlers CreateOrder
func (s *serverAPI) CreateOrder(
	ctx context.Context,
	req *orderv1.CreateOrderRequest,
) (*orderv1.CreateOrderResponse, error) {
	err := validateItemAndQuan(req.Item, req.Quantity)
	if err != nil {
		return nil, err
	}

	// Servise layer
	id, err := s.crud.Create(ctx, req.Item, req.Quantity)
	if err != nil {
		slog.Debug(
			"error in servise layer",
			slog.String("hendler: ", "CreateOrder"),
			slog.String("Error", err.Error()),
		)

		return nil, err
	}

	return &orderv1.CreateOrderResponse{Id: id}, nil
}

// Hanlers GetOrder
func (s *serverAPI) GetOrder(
	ctx context.Context,
	req *orderv1.GetOrderRequest,
) (*orderv1.GetOrderResponse, error) {
	err := validateID(req.Id)
	if err != nil {
		return nil, err
	}

	// Servise layer
	order, err := s.crud.Get(ctx, req.Id)
	if err != nil {
		slog.Debug(
			"error in servise layer",
			slog.String("hendler: ", "GetOrder"),
			slog.String("Error", err.Error()),
		)

		return nil, err
	}

	return &orderv1.GetOrderResponse{Order: conversionTypeOrders(order)}, nil
}

// Hanlers UpdateOrder
func (s *serverAPI) UpdateOrder(
	ctx context.Context,
	req *orderv1.UpdateOrderRequest,
) (*orderv1.UpdateOrderResponse, error) {
	err := validateID(req.Id)
	if err != nil {
		return nil, err
	}

	err = validateItemAndQuan(req.Item, req.Quantity)
	if err != nil {
		return nil, err
	}

	// Servise layer
	order, err := s.crud.Update(ctx, req.Id, req.Item, req.Quantity)
	if err != nil {
		slog.Debug(
			"error in servise layer",
			slog.String("hendler: ", "UpdateOrder"),
			slog.String("Error", err.Error()),
		)

		return nil, err
	}

	return &orderv1.UpdateOrderResponse{Order: conversionTypeOrders(order)}, nil
}

// Hanlers DeleteOrder
func (s *serverAPI) DeleteOrder(
	ctx context.Context,
	req *orderv1.DeleteOrderRequest,
) (*orderv1.DeleteOrderResponse, error) {
	err := validateID(req.Id)
	if err != nil {
		return nil, err
	}

	// Servise layer
	success, err := s.crud.Delete(ctx, req.Id)
	if err != nil {
		slog.Debug(
			"error in servise layer",
			slog.String("hendler: ", "DeleteOrder"),
			slog.String("Error", err.Error()),
		)

		return nil, err
	}

	return &orderv1.DeleteOrderResponse{Success: success}, nil
}

// Hanlers ListOrder
func (s *serverAPI) ListOrders(
	ctx context.Context,
	req *orderv1.ListOrdersRequest,
) (*orderv1.ListOrdersResponse, error) {
	// Servise layer
	orders, err := s.crud.List(ctx)
	if err != nil {
		slog.Debug(
			"error in servise layer",
			slog.String("hendler: ", "ListOrders"),
			slog.String("Error", err.Error()),
		)

		return nil, err
	}

	grpcOrder := make([]*orderv1.Order, 0, len(orders))
	for _, order := range orders {
		grpcOrder = append(grpcOrder, conversionTypeOrders(order))
	}

	return &orderv1.ListOrdersResponse{Orders: grpcOrder}, nil
}

// Validate id
func validateID(id string) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "id invalid")
	}

	return nil
}

// Validate item and quantity
func validateItemAndQuan(item string, quan int32) error {
	if item == "" {
		return status.Error(codes.InvalidArgument, "invalid item")
	}

	if quan < 0 {
		return status.Error(codes.InvalidArgument, " quantity is less than zero")
	}

	return nil
}

func conversionTypeOrders(order *models.Order) *orderv1.Order {
	return &orderv1.Order{
		Id:       order.ID,
		Item:     order.Item,
		Quantity: order.Quantity,
	}
}
