package tests

import (
	"order/internal/domain/models"
	"order/tests/client"
	"testing"

	orderv1 "github.com/Nikitosnim/protos/gen/go/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestDeleteOrder_Success(t *testing.T) {
	ctx, st := client.New(t)
	testOrders := make([]models.Order, 0, 5)

	for range 400 {
		newOrder := newOrderRandom(t)

		res, err := st.CrudClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
			Item:     newOrder.Item,
			Quantity: newOrder.Quantity,
		})

		require.NoError(t, err, err)
		assert.NotEmpty(t, res.Id)

		newOrder.ID = res.Id
		testOrders = append(testOrders, newOrder)
	}

	// res, _ := st.CrudClient.ListOrders(ctx, &orderv1.ListOrdersRequest{})
	// testOrders := res.Orders

	const connectMax = 20
	sem := make(chan struct{}, connectMax)
	g, gctx := errgroup.WithContext(ctx)

	for _, order := range testOrders {
		sem <- struct{}{}
		g.Go(func() error {
			defer func() { <-sem }()

			res, err := st.CrudClient.DeleteOrder(gctx, &orderv1.DeleteOrderRequest{
				Id: order.ID,
			})
			if err != nil {
				return err
			}
			if !res.Success {
				return assert.AnError
			}
			return nil
		})
	}
	require.NoError(t, g.Wait())
}
