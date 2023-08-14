// Package repository contains getting and putting info out of the project for example work with db or grpc stream
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// TradingServiceRepository contains an object of *pgxpool.Pool
type TradingServiceRepository struct {
	pool *pgxpool.Pool
}

// NewTradingServiceRepository is the constructor for TradingServiceRepository
func NewTradingServiceRepository(pool *pgxpool.Pool) *TradingServiceRepository {
	return &TradingServiceRepository{pool: pool}
}

// AddPosition adds info about opened position to database
func (r *TradingServiceRepository) AddPosition(ctx context.Context, position *model.Position, shareStartPrice float64) error {
	shareStartPriceDecimal := decimal.NewFromFloat(shareStartPrice)
	amountMoneyDecimal := decimal.NewFromFloat(position.MoneyAmount)
	shareAmountDecimal := amountMoneyDecimal.Div(shareStartPriceDecimal)
	shareAmount := shareAmountDecimal.InexactFloat64()
	_, err := r.pool.Exec(ctx, `INSERT INTO positions (positionid, profileid, vector, shareName, shareAmount, 
		shareStartPrice, stopLoss, takeProfit, openedTime) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		position.PositionID, position.ProfileID, position.Vector, position.ShareName, shareAmount, shareStartPrice, position.StopLoss, position.TakeProfit, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("TradingServiceRepository -> AddPosition -> %w", err)
	}
	return nil
}

// ClosePosition adds info about closed position to database
func (r *TradingServiceRepository) ClosePosition(ctx context.Context, positionID uuid.UUID, shareEndPrice float64) error {
	res, err := r.pool.Exec(ctx, "UPDATE positions SET shareEndPrice = $1, closedTime = $2 WHERE positionid = $3", shareEndPrice, time.Now().UTC(), positionID)
	if err != nil {
		return fmt.Errorf("TradingServiceRepository -> ClosePosition -> %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("TradingServiceRepository -> ClosePosition -> %w", pgx.ErrNoRows)
	}
	return nil
}
