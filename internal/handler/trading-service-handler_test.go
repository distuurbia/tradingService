package handler

import (
	"context"
	"testing"

	"github.com/distuurbia/tradingService/internal/handler/mocks"
	"github.com/distuurbia/tradingService/internal/model"
	"github.com/distuurbia/tradingService/protocol/trading"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidationID(t *testing.T) {
	h := NewTradingServiceHandler(nil, validate, nil)
	testID := uuid.New()
	validatedID, err := h.validationID(context.Background(), testID.String())
	require.NoError(t, err)
	require.Equal(t, testID, validatedID)
}

func TestAddPosition(t *testing.T) {
	s := new(mocks.TradingServiceService)
	s.On("AddPosition", mock.Anything, mock.Anything).Return(nil)
	h := NewTradingServiceHandler(s, validate, &cfg)
	_, err := h.AddPosition(context.Background(), &trading.AddPositionRequest{
		Position: testProtoPosition,
	})
	require.NoError(t, err)
}

func TestClosePosition(t *testing.T) {
	testPnl := 50.247
	s := new(mocks.TradingServiceService)
	s.On("ClosePosition", mock.Anything, mock.Anything, mock.Anything).Return(testPnl, nil)

	h := NewTradingServiceHandler(s, validate, &cfg)
	resp, err := h.ClosePosition(context.Background(), &trading.ClosePositionRequest{
		ProfileID:  uuid.New().String(),
		PositionID: uuid.New().String(),
	})
	require.NoError(t, err)
	require.Equal(t, testPnl, resp.Pnl)
}

func TestReadAllOpenedPositionsByProfileID(t *testing.T) {
	openedPositions := make([]*model.OpenedPosition, 0)
	openedPositions = append(openedPositions, &testOpenedPosition, &testOpenedPosition)
	s := new(mocks.TradingServiceService)
	s.On("ReadAllOpenedPositionsByProfileID", mock.Anything, mock.Anything).Return(openedPositions, nil)

	h := NewTradingServiceHandler(s, validate, nil)
	resp, err := h.ReadAllOpenedPositionsByProfileID(context.Background(), &trading.ReadAllOpenedPositionsByProfileIDRequest{
		ProfileID: uuid.New().String(),
	})
	require.NoError(t, err)
	require.Equal(t, len(openedPositions), len(resp.OpenedPositions))
}
