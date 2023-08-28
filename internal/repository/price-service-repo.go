// Package repository contains getting and putting info out of the project for example work with db or grpc stream
package repository

import (
	"context"
	"fmt"
	"strings"

	priceProtocol "github.com/distuurbia/PriceService/protocol/price"
	"github.com/distuurbia/tradingService/internal/config"
	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
)

// PriceServiceRepository contains an object of priceProtocol.PriceServiceServiceClient
type PriceServiceRepository struct {
	client priceProtocol.PriceServiceServiceClient
	cfg    *config.Config
}

// NewPriceServiceRepository is the constructor for PriceServiceRepository
func NewPriceServiceRepository(client priceProtocol.PriceServiceServiceClient, cfg *config.Config) *PriceServiceRepository {
	return &PriceServiceRepository{client: client, cfg: cfg}
}

// Subscribe receives the data from grpc stream and fills the given channel
func (r *PriceServiceRepository) Subscribe(ctx context.Context, subID uuid.UUID, subscribersShares chan model.Share, errSubscribe chan error) {
	selectedShares := strings.Split(r.cfg.TradingServiceShares, ",")
	stream, err := r.client.Subscribe(ctx, &priceProtocol.SubscribeRequest{
		SelectedShares: selectedShares,
		UUID:           subID.String()})
	if err != nil {
		errSubscribe <- fmt.Errorf("PriceServiceRepository -> Subscribe -> %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			protoResponse, err := stream.Recv()
			if err != nil {
				errSubscribe <- fmt.Errorf("PriceServiceRepository -> Subscribe -> %w", err)
				return
			}
			errSubscribe <- nil

			for _, protoShare := range protoResponse.Shares {
				subscribersShares <- model.Share{Name: protoShare.Name, Price: protoShare.Price}
			}
		}
	}
}
