package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/caarlos0/env"
	priceProtocol "github.com/distuurbia/PriceService/protocol/price"
	"github.com/distuurbia/PriceService/protocol/price/mocks"
	"github.com/distuurbia/tradingService/internal/config"
	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockClient := new(mocks.PriceServiceServiceClient)

	mockStream := new(mocks.PriceServiceService_SubscribeClient)

	expectedShares := []*priceProtocol.Share{
		{Name: "Apple", Price: 100},
		{Name: "Tesla", Price: 200},
	}
	mockStream.On("Recv").Return(&priceProtocol.SubscribeResponse{Shares: expectedShares}, nil)

	mockClient.On("Subscribe", mock.Anything, mock.Anything).Return(mockStream, nil)

	var cfg config.Config
	err := env.Parse(&cfg)
	require.NoError(t, err)

	lenOfReadableShares := len(strings.Split(cfg.TradingServiceShares, ","))

	r := NewPriceServiceRepository(mockClient, &cfg)

	subID := uuid.New()
	subscribersShares := make(chan model.Share, lenOfReadableShares)
	errSubscribe := make(chan error)

	go r.Subscribe(ctx, subID, subscribersShares, errSubscribe)

	select {
	case err := <-errSubscribe:
		require.NoError(t, err)
	case <-time.After(time.Second * 5):
		t.Error("Timed out while waiting for error")
	}

	select {
	case firstShare := <-subscribersShares:
		require.Equal(t, len(expectedShares), len(subscribersShares)+1)

		require.Equal(t, firstShare.Name, expectedShares[0].Name)
		require.Equal(t, firstShare.Price, expectedShares[0].Price)

		secondShare := <-subscribersShares
		require.Equal(t, secondShare.Name, expectedShares[1].Name)
		require.Equal(t, secondShare.Price, expectedShares[1].Price)

	case <-time.After(time.Second * 5):
		t.Error("Timed out while waiting for shares")
	}

	mockClient.AssertExpectations(t)
	mockStream.AssertExpectations(t)
}
