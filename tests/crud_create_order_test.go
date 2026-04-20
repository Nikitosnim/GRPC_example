package tests

import (
	"crypto/rand"
	"encoding/base64"
	"order/internal/domain/models"
	"order/tests/client"
	"testing"

	orderv1 "github.com/Nikitosnim/protos/gen/go/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestCreateOrder_Success(t *testing.T) {
	testOrders := make([]models.Order, 0, 5)
	for range 40 {
		newOrder := newOrderRandom(t)
		testOrders = append(testOrders, newOrder)
	}

	ctx, st := client.New(t)

	const connectMax = 20
	sem := make(chan struct{}, connectMax)
	g, gctx := errgroup.WithContext(ctx)

	for _, order := range testOrders {
		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }()

			res, err := st.CrudClient.CreateOrder(gctx, &orderv1.CreateOrderRequest{
				Item:     order.Item,
				Quantity: order.Quantity,
			})
			if err != nil {
				return err
			}
			if res.Id == "" {
				return assert.AnError
			}

			return nil
		})

	}

	require.NoError(t, g.Wait())
}

func newOrderRandom(t *testing.T) models.Order {
	bufId := make([]byte, 10)
	_, err := rand.Read(bufId)
	require.NoError(t, err)

	bufItem := make([]byte, 10)
	_, err = rand.Read(bufItem)
	require.NoError(t, err)

	prime, err := rand.Prime(rand.Reader, 30)
	require.NoError(t, err)

	newId := base64.URLEncoding.EncodeToString(bufId)[:10]
	newItem := base64.URLEncoding.EncodeToString(bufItem)[:10]
	newQuan := int32(prime.Int64())

	return models.Order{
		ID:       newId,
		Item:     newItem,
		Quantity: newQuan,
	}
}
