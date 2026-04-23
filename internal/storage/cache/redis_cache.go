package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"order/internal/config"
	"order/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	db  *redis.Client
	cfg *config.Config
}

func New(ctx context.Context, cfg *config.Config) (*cache, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Cache.Addr,
		Password:     cfg.Cache.Password,
		DB:           cfg.Cache.DB,
		Username:     cfg.Cache.User,
		MaxRetries:   cfg.Cache.MaxRetries,
		DialTimeout:  cfg.Cache.DialTimeout,
		ReadTimeout:  cfg.Cache.Timeout,
		WriteTimeout: cfg.Cache.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		db.Close()
		return nil, err
	}

	return &cache{
		db:  db,
		cfg: cfg,
	}, nil
}

func (c *cache) Get(ctx context.Context, id string) (*models.Order, error) {
	key := "order:" + id
	val, err := c.db.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	slog.Info(
		"cache",
		slog.String("val", val),
	)

	var order models.Order
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *cache) Set(ctx context.Context, order *models.Order) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	key := "order:" + order.ID
	err = c.db.Set(ctx, key, data, c.cfg.Cache.TTL).Err()
	return err
}

func (c *cache) Delete(ctx context.Context, id string)  error {
	key := "order:" + id
	err := c.db.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) Stop() {
	c.db.Close()
}
