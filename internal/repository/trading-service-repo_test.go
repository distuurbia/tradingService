package repository

import (
	"context"
	"testing"
	"time"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBackupAllOpenedPositions(t *testing.T) {
	testPosition.ProfileID = uuid.New()
	testPosition.PositionID = uuid.New()
	shareStartPrice := float64(25)
	shareAmount := float64(3)

	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	testPosition.ProfileID = uuid.New()
	testPosition.PositionID = uuid.New()

	err = r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	openedPositions, err := r.BackupAllOpenedPositions(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, len(openedPositions))
}

func TestAddPosition(t *testing.T) {
	shareStartPrice := float64(25)
	shareAmount := float64(3)
	testPosition.ProfileID = uuid.New()
	testPosition.PositionID = uuid.New()
	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	var readPosition model.Position
	var readStartPrice float64
	var readShareAmount float64
	var openedTime time.Time

	err = r.pool.QueryRow(context.Background(), `SELECT profileid, shortOrLong, sharename, sharestartprice, shareamount, stoploss, takeprofit, openedTime
	 	FROM positions WHERE positionid = $1`, testPosition.PositionID).Scan(&readPosition.ProfileID, &readPosition.ShortOrLong, &readPosition.ShareName, &readStartPrice,
		&readShareAmount, &readPosition.StopLoss, &readPosition.TakeProfit, &openedTime)

	require.NoError(t, err)

	require.Equal(t, testPosition.ProfileID, readPosition.ProfileID)
	require.Equal(t, testPosition.ShortOrLong, readPosition.ShortOrLong)
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

func TestReadAllOpenedPositionsByProfileID(t *testing.T) {
	testPosition.ProfileID = uuid.New()
	testPosition.PositionID = uuid.New()
	shareStartPrice := float64(25)
	shareAmount := float64(3)

	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	testPosition.PositionID = uuid.New()

	err = r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	openedPositions, err := r.ReadAllOpenedPositionsByProfileID(context.Background(), testPosition.ProfileID)
	require.NoError(t, err)
	require.Equal(t, 2, len(openedPositions))
}

func TestReadPositionRow(t *testing.T) {
	shareStartPrice := float64(25)
	shareEndPrice := float64(30)
	shareAmount := float64(3)
	testPosition.PositionID = uuid.New()

	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)
	err = r.ClosePosition(context.Background(), testPosition.PositionID, shareEndPrice)
	require.NoError(t, err)

	readShareStartPrice, readShareEndPrice, readShareAmount, readShortOrLong, err := r.ReadPositionRow(context.Background(), testPosition.PositionID)
	require.NoError(t, err)
	require.Equal(t, shareStartPrice, readShareStartPrice)
	require.Equal(t, shareEndPrice, readShareEndPrice)
	require.Equal(t, shareAmount, readShareAmount)
	require.Equal(t, testPosition.ShortOrLong, readShortOrLong)
}

func TestReadShareNameByPositionID(t *testing.T) {
	shareStartPrice := float64(25)
	shareAmount := float64(3)
	testPosition.PositionID = uuid.New()

	err := r.AddPosition(context.Background(), &testPosition, shareAmount, shareStartPrice)
	require.NoError(t, err)

	shareName, err := r.ReadShareNameByPositionID(context.Background(), testPosition.PositionID)
	require.NoError(t, err)
	require.Equal(t, testPosition.ShareName, shareName)
}
