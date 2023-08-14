// Package repository contains getting and putting info out of the project for example work with db or grpc stream
package repository

import (
	"context"
	"fmt"

	balanceProtocol "github.com/distuurbia/balance/protocol/balance"
	"github.com/shopspring/decimal"
)

// BalanceRepository contains an object of balanceProtocol.BalanceServiceClient
type BalanceRepository struct {
	client balanceProtocol.BalanceServiceClient
}

// NewBalanceRepository is the constructor for BalanceRepository
func NewBalanceRepository(client balanceProtocol.BalanceServiceClient) *BalanceRepository {
	return &BalanceRepository{client: client}
}

// WithdrawOnPosition takes money from balance of exact profile until the position is closed
func (r *BalanceRepository) WithdrawOnPosition(ctx context.Context, profileID string, amount float64) error {
	getBalanceResponse, err := r.client.GetBalance(ctx, &balanceProtocol.GetBalanceRequest{ProfileID: profileID})
	if err != nil {
		return fmt.Errorf("BalanceRepository -> WithdrawOnPosition -> %w", err)
	}
	if getBalanceResponse.TotalBalance < amount {
		return fmt.Errorf("BalanceRepository -> WithdrawOnPosition -> error: not enough money on your balance")
	}
	coef := decimal.NewFromFloat(-1)
	decimalAmount := decimal.NewFromFloat(amount)
	decimalAmount = decimalAmount.Mul(coef)
	_, err = r.client.AddBalanceChange(ctx, &balanceProtocol.AddBalanceChangeRequest{ProfileID: profileID, Amount: decimalAmount.InexactFloat64()})
	if err != nil {
		return fmt.Errorf("BalanceRepository -> WithdrawOnPosition -> %w", err)
	}
	return nil
}

// MoneyBackWithPnl returnes money to balance of exact profile when the position is closed
func (r *BalanceRepository) MoneyBackWithPnl(ctx context.Context, profileID string, moneyBack float64) error {
	_, err := r.client.AddBalanceChange(ctx, &balanceProtocol.AddBalanceChangeRequest{ProfileID: profileID, Amount: moneyBack})
	if err != nil {
		return fmt.Errorf("BalanceRepository -> WithdrawOnPosition -> %w", err)
	}
	return nil
}
