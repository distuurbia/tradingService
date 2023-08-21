package repository

import (
	"context"
	"testing"
	"time"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAddPosition(t *testing.T) {
	shareStartPrice := float64(25)
	shareAmount := float64(3)

	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	var readPosition model.Position
	var readStartPrice float64
	var readShareAmount float64
	var openedTime time.Time

	err = r.pool.QueryRow(context.Background(), `SELECT profileid, vector, sharename, sharestartprice, shareamount, stoploss, takeprofit, openedTime
	 	FROM positions WHERE positionid = $1`, testPosition.PositionID).Scan(&readPosition.ProfileID, &readPosition.Vector, &readPosition.ShareName, &readStartPrice,
		&readShareAmount, &readPosition.StopLoss, &readPosition.TakeProfit, &openedTime)

	require.NoError(t, err)

	require.Equal(t, testPosition.ProfileID, readPosition.ProfileID)
	require.Equal(t, testPosition.Vector, readPosition.Vector)
	require.Equal(t, testPosition.ShareName, readPosition.ShareName)
	require.Equal(t, shareStartPrice, readStartPrice)
	require.Equal(t, shareAmount, readShareAmount)
	require.Equal(t, testPosition.StopLoss, readPosition.StopLoss)
	require.Equal(t, testPosition.TakeProfit, readPosition.TakeProfit)
	require.NotEmpty(t, openedTime)
}

func TestClosePosition(t *testing.T) {
	shareStartPrice := float64(25)
	shareEndPrice := float64(30)
	shareAmount := float64(3)
	testPosition.PositionID = uuid.New()

	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)
	err = r.ClosePosition(context.Background(), testPosition.PositionID, shareEndPrice)
	require.NoError(t, err)

	var readShareEndPrice float64
	var readClosedTime time.Time
	err = r.pool.QueryRow(context.Background(), "SELECT shareendprice, closedtime FROM positions WHERE positionid = $1", testPosition.PositionID).
		Scan(&readShareEndPrice, &readClosedTime)

	require.NoError(t, err)
	require.Equal(t, shareEndPrice, readShareEndPrice)
	require.NotEmpty(t, readClosedTime)
}

// func TestReadShareNameByID(t *testing.T) {
// 	shareStartPrice := float64(25)
// 	testPosition.PositionID = uuid.New()

// 	err := r.AddPosition(context.Background(), &testPosition, shareStartPrice)
// 	require.NoError(t, err)

// 	readShareName, err := r.ReadShareNameByID(context.Background(), testPosition.PositionID)

// 	require.NoError(t, err)
// 	require.Equal(t, testPosition.ShareName, readShareName)
// }
