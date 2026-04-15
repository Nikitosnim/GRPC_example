package order

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"order/internal/domain/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Order struct {
	storage Storage
}

// Interfase order storage
type Storage interface {
	Create(
		ctx context.Context,
		id string,
		item string,
		quan int32,
	) error

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

	Stop()
}

// New rerturns a new instance of the Order service
func New(storage Storage) *Order {
	return &Order{
		storage: storage,
	}
}

func (o *Order) Create(
	ctx context.Context,
	item string,
	quan int32,
) (string, error) {
	var id string

	buf := make([]byte, 15)
	_, err := rand.Read(buf)
	if err != nil {
		return "", status.Error(codes.Internal, "internal err")
	}

	id = base64.URLEncoding.EncodeToString(buf)[:15]
	err = o.storage.Create(ctx, id, item, quan)
	if err != nil {
		slog.Debug(
			"error in storage layer",
			slog.String("Serv: ", "Create"),
			slog.String("Error", err.Error()),
		)

		return "", status.Error(codes.Internal, "internal err")
	}

	const op = "services.Create"
	slog.Info(
		"method Create",
		slog.String("op", op),
		slog.String("id", id),
	)

	return id, nil
}

func (o *Order) Get(
	ctx context.Context,
	id string,
) (*models.Order, error) {
	const op = "services.Get"
	slog.Info(
		"method Get",
		slog.String("op", op),
		slog.String("id", id),
	)

	order, err := o.storage.Get(ctx, id)
	if err != nil {
		slog.Debug(
			"error in storage layer",
			slog.String("Serv: ", "Get"),
			slog.String("Error", err.Error()),
		)

		return nil, status.Error(codes.Internal, "internal err")
	}

	return order, nil
}

func (o *Order) Update(
	ctx context.Context,
	id string,
	item string,
	quan int32,
) (*models.Order, error) {
	const op = "services.Update"
	slog.Info(
		"method Update",
		slog.String("op", op),
		slog.String("id", id),
	)

	order, err := o.storage.Update(ctx, id, item, quan)
	if err != nil {
		slog.Debug(
			"error in storage layer",
			slog.String("Serv: ", "Update"),
			slog.String("Error", err.Error()),
		)

		return nil, status.Error(codes.Internal, "internal err")
	}

	return order, nil
}

func (o *Order) Delete(
	ctx context.Context,
	id string,
) (bool, error) {
	const op = "storage.Delete"
	slog.Info(
		"method Delete",
		slog.String("op", op),
		slog.String("id", id),
	)

	success, err := o.storage.Delete(ctx, id)
	if err != nil {
		slog.Debug(
			"error in storage layer",
			slog.String("Serv: ", "Delete"),
			slog.String("Error", err.Error()),
		)

		return success, status.Error(codes.Internal, "internal err")
	}

	return success, nil
}

func (o *Order) List(ctx context.Context) ([]*models.Order, error) {
	const op = "services.List"
	slog.Info(
		"method List",
		slog.String("op", op),
	)

	orders, err := o.storage.List(ctx)
	if err != nil {
		slog.Debug(
			"error in storage layer",
			slog.String("Serv: ", "List"),
			slog.String("Error", err.Error()),
		)

		return nil, status.Error(codes.Internal, "internal err")
	}

	return orders, nil
}
