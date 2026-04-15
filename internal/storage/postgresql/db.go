package postgresql

import (
	"context"
	"log/slog"
	"order/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

const (
	createQuery = `
	INSERT INTO orders (id, item, quantity)
	VALUES ($1, $2, $3)
	`

	getQuery = `SELECT item, quantity FROM orders WHERE id = $1`

	updateQuery = `
	UPDATE orders 
	SET item = $2, quantity = $3 
	WHERE id = $1 
	RETURNING id, item, quantity
	`

	deleteQuery = `DELETE FROM orders WHERE id = $1`

	listQuery = `SELECT id, item, quantity FROM orders`
)

func New(ctx context.Context, coonString string) (*Storage, error) {
	conf, err := pgxpool.ParseConfig(coonString)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = 10
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 1 * time.Minute
	conf.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) Create(
	ctx context.Context,
	id string,
	item string,
	quan int32,
) error {
	const op = "storage.Create"
	slog.Info(
		"method Create",
		slog.String("op", op),
		slog.String("id", id),
	)

	_, err := s.pool.Exec(ctx, createQuery, id, item, quan)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get(
	ctx context.Context,
	id string,
) (*models.Order, error) {
	var (
		item string
		quan int32
	)

	const op = "storage.Get"
	slog.Info(
		"method Get",
		slog.String("op", op),
		slog.String("id", id),
	)

	err := s.pool.QueryRow(ctx, getQuery, id).Scan(&item, &quan)
	if err != nil {
		return nil, err
	}

	return &models.Order{
		ID:       id,
		Item:     item,
		Quantity: quan,
	}, nil
}

func (s *Storage) Update(
	ctx context.Context,
	id string,
	item string,
	quan int32,
) (*models.Order, error) {
	var updateOrder models.Order

	const op = "storage.Update"
	slog.Info(
		"method Update",
		slog.String("op", op),
		slog.String("id", id),
	)

	err := s.pool.QueryRow(ctx, updateQuery, id, item, quan).Scan(
		&updateOrder.ID,
		&updateOrder.Item,
		&updateOrder.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &updateOrder, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	id string,
) (bool, error) {
	const op = "storage.Delete"
	slog.Info(
		"method Delete",
		slog.String("op", op),
		slog.String("id", id),
	)

	res, err := s.pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, nil
}

func (s *Storage) List(ctx context.Context) ([]*models.Order, error) {
	const op = "storage.List"
	slog.Info(
		"method List",
		slog.String("op", op),
	)

	rows, err := s.pool.Query(ctx, listQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sliseOrders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.Item, &order.Quantity)
		if err != nil {
			return nil, err
		}

		sliseOrders = append(sliseOrders, &models.Order{
			ID:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sliseOrders, nil
}

func (s *Storage) Stop() {
	s.pool.Close()
}
