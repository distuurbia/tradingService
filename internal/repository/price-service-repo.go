// Package repository contains getting and putting info out of the project for example work with db or grpc stream
package repository

import (
	"context"
	"fmt"

	priceProtocol "github.com/distuurbia/PriceService/protocol/price"
	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
)

// PriceServiceRepository contains an object of priceProtocol.PriceServiceServiceClient
type PriceServiceRepository struct {
	client priceProtocol.PriceServiceServiceClient
}

// NewPriceServiceRepository is the constructor for PriceServiceRepository
func NewPriceServiceRepository(client priceProtocol.PriceServiceServiceClient) *PriceServiceRepository {
	return &PriceServiceRepository{client: client}
}

// Subscribe receive data from grpc stream and fills the given channel
func (r *PriceServiceRepository) Subscribe(ctx context.Context, subID uuid.UUID, subscribersShares chan []*model.Share, errSubscribe chan error) {
	selectedShares := make([]string, 0)
	selectedShares = append(selectedShares, "Apple", "Twitter", "Tesla", "Coca-Cola", "Berkshire Hathaway Inc.")

	stream, err := r.client.Subscribe(ctx, &priceProtocol.SubscribeRequest{
		SelectedShares: selectedShares,
		UUID:           subID.String()})
	if err != nil {
		errSubscribe <- fmt.Errorf("PriceServiceRepository -> Subscribe -> %w", err)
	}

	for {
		protoResponse, err := stream.Recv()

		if err != nil {
			errSubscribe <- fmt.Errorf("PriceServiceRepository -> Subscribe -> %w", err)
		}
		recievedShares := make([]*model.Share, 0)

		for _, protoShare := range protoResponse.Shares {
			recievedShares = append(recievedShares, &model.Share{Name: protoShare.Name, Price: protoShare.Price})
		}

		errSubscribe <- nil
		subscribersShares <- recievedShares
	}
}
