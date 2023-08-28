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
func (r *TradingServiceRepository) AddPosition(ctx context.Context, position *model.Position, shareAmount, shareStartPrice float64) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO positions (positionid, profileid, shortOrLong, shareName, shareAmount, 
		shareStartPrice, stopLoss, takeProfit, openedTime) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		position.PositionID, position.ProfileID, position.ShortOrLong, position.ShareName, shareAmount, shareStartPrice, position.StopLoss,
		position.TakeProfit, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("TradingServiceRepository -> AddPosition -> %w", err)
	}
	return nil
}

// ClosePosition adds shareEndPrice and closedTime values to exact position row using positionID
func (r *TradingServiceRepository) ClosePosition(ctx context.Context, positionID uuid.UUID, shareEndPrice float64) error {
	res, err := r.pool.Exec(ctx, "UPDATE positions SET shareEndPrice = $1, closedTime = $2 WHERE positionid = $3", shareEndPrice,
		time.Now().UTC(), positionID)
	if err != nil {
		return fmt.Errorf("TradingServiceRepository -> ClosePosition -> %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("TradingServiceRepository -> ClosePosition -> %w", pgx.ErrNoRows)
	}
	return nil
}

// ReadAllOpenedPositionsByProfileID reads from db all positions that wasn't be closed for an exact profile using its ID
func (r *TradingServiceRepository) ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error) {
	openedPositions := make([]*model.OpenedPosition, 0)

	rows, err := r.pool.Query(ctx, `SELECT  profileID, positionID, shortOrLong, shareName, shareAmount, 
		shareStartPrice, stopLoss, takeProfit, openedTime FROM positions WHERE profileID = $1 AND closedTime is NULL`, profileID)
	if err != nil {
		return nil, fmt.Errorf("TradingServiceRepository -> BackupAllOpenedPositions -> %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		readedPosition := model.OpenedPosition{}
		err := rows.Scan(&readedPosition.ProfileID, &readedPosition.PositionID, &readedPosition.ShortOrLong, &readedPosition.ShareName,
			&readedPosition.ShareAmount, &readedPosition.ShareStartPrice, &readedPosition.StopLoss, &readedPosition.TakeProfit,
			&readedPosition.OpenedTime)
		if err != nil {
			return nil, fmt.Errorf("TradingServiceRepository -> BackupAllOpenedPositions -> %w", err)
		}
		shareAmountDecimal := decimal.NewFromFloat(readedPosition.ShareAmount)
		shareStartPriceDecimal := decimal.NewFromFloat(readedPosition.ShareStartPrice)
		moneyAmountDecimal := shareAmountDecimal.Mul(shareStartPriceDecimal)
		readedPosition.MoneyAmount = moneyAmountDecimal.InexactFloat64()

		openedPositions = append(openedPositions, &readedPosition)
	}

	return openedPositions, nil
}

// BackupAllOpenedPositions reads from db all positions that wasn't be closed to backup it
func (r *TradingServiceRepository) BackupAllOpenedPositions(ctx context.Context) ([]*model.OpenedPosition, error) {
	openedPositions := make([]*model.OpenedPosition, 0)

	rows, err := r.pool.Query(ctx, `SELECT  profileID, positionID, shortOrLong, shareName, shareAmount, 
		shareStartPrice, stopLoss, takeProfit FROM positions WHERE closedTime is NULL`)
	if err != nil {
		return nil, fmt.Errorf("TradingServiceRepository -> BackupAllOpenedPositions -> %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		readedPosition := model.OpenedPosition{}
		err := rows.Scan(&readedPosition.ProfileID, &readedPosition.PositionID, &readedPosition.ShortOrLong, &readedPosition.ShareName,
			&readedPosition.ShareAmount, &readedPosition.ShareStartPrice, &readedPosition.StopLoss, &readedPosition.TakeProfit)
		if err != nil {
			return nil, fmt.Errorf("TradingServiceRepository -> BackupAllOpenedPositions -> %w", err)
		}
		shareAmountDecimal := decimal.NewFromFloat(readedPosition.ShareAmount)
		shareStartPriceDecimal := decimal.NewFromFloat(readedPosition.ShareStartPrice)
		moneyAmountDecimal := shareAmountDecimal.Mul(shareStartPriceDecimal)
		readedPosition.MoneyAmount = moneyAmountDecimal.InexactFloat64()

		openedPositions = append(openedPositions, &readedPosition)
	}

	return openedPositions, nil
}

// ReadPositionRow returnes share start price and share amount values for an exact position using positionID
func (r *TradingServiceRepository) ReadPositionRow(ctx context.Context, positionID uuid.UUID) (shareStartPrice, shareEndPrice,
	shareAmount float64, shortOrLong string, err error) {
	err = r.pool.QueryRow(ctx, "SELECT shareStartPrice, shareEndPrice, shareAmount, shortOrLong FROM positions WHERE positionID = $1",
		positionID).Scan(&shareStartPrice, &shareEndPrice, &shareAmount, &shortOrLong)

	if err != nil {
		return 0, 0, 0, "", fmt.Errorf("TradingServiceRepository -> ReadInfoAboutOpenedPosition -> %w", err)
	}

	return shareStartPrice, shareEndPrice, shareAmount, shortOrLong, nil
}

// ReadShareNameByPositionID returnes share name value for an exact position by its ID
func (r *TradingServiceRepository) ReadShareNameByPositionID(ctx context.Context, positionID uuid.UUID) (string, error) {
	var shareName string
	err := r.pool.QueryRow(ctx, "SELECT shareName FROM positions WHERE positionID = $1",
		positionID).Scan(&shareName)

	if err != nil {
		return "", fmt.Errorf("TradingServiceRepository -> ReadShareNameByPositionID -> %w", err)
	}

	return shareName, nil
}
