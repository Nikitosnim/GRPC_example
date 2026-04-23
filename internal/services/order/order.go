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
	cache   Cache
	storage Storage
}

type Cache interface {
	Get(ctx context.Context, id string) (*models.Order, error)
	Set(ctx context.Context, order *models.Order) error
	Delete(ctx context.Context, id string)  error
	Stop()
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
func New(storage Storage, cache Cache) *Order {
	return &Order{
		storage: storage,
		cache:   cache,
	}
}

func (o *Order) Create(
	ctx context.Context,
	item string,
	quan int32,
) (string, error) {
	const op = "services.Create"

	id, err := GenerateId()
	if err != nil {
		return "", err
	}

	err = o.storage.Create(ctx, id, item, quan)
	if err != nil {
		slog.Debug(
			"error in storage layer",
			slog.String("Serv: ", "Create"),
			slog.String("Error", err.Error()),
		)

		return "", status.Error(codes.Internal, "internal err")
	}

	err = o.cache.Set(ctx, &models.Order{
		ID:       id,
		Item:     item,
		Quantity: quan,
	})
	if err != nil {
		slog.Debug(
			"error in cache",
			slog.String("Serv: ", "Create"),
			slog.String("Error", err.Error()),
		)
	}

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

	order, err := o.cache.Get(ctx, id)
	if err != nil {
		slog.Debug(
			"error in cache",
			slog.String("Serv: ", "Get"),
			slog.String("Error", err.Error()),
		)
	}

	if order == nil {
		order, err = o.storage.Get(ctx, id)
		if err != nil {
			slog.Debug(
				"error in storage layer",
				slog.String("Serv: ", "Get"),
				slog.String("Error", err.Error()),
			)

			return nil, status.Error(codes.Internal, "internal err")
		}

		o.cache.Set(ctx, order)
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

	_ = o.cache.Delete(ctx, id)
	err = o.cache.Set(ctx, &models.Order{
		ID: id,
		Item: item,
		Quantity: quan,
	})
	
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

	_ = o.cache.Delete(ctx, id)

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

func GenerateId() (string, error) {
	buf := make([]byte, 15)
	_, err := rand.Read(buf)
	if err != nil {
		return "", status.Error(codes.Internal, "internal err")
	}

	return base64.URLEncoding.EncodeToString(buf)[:15], nil
}
