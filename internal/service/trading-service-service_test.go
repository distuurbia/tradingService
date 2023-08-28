package service

import (
	"context"
	"testing"
	"time"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/distuurbia/tradingService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendSharesToProfiles(t *testing.T) {
	priceRps := new(PriceServiceRepositoryHelper)

	s := NewTradingServiceService(priceRps, nil, nil)
	go s.SendSharesToProfiles(context.Background(), 2)
	time.Sleep(time.Millisecond)
	require.Equal(t, testShareApple.Price, s.profMngr.Shares[testShareApple.Name])
	require.Equal(t, testShareTesla.Price, s.profMngr.Shares[testShareTesla.Name])
}

func TestOpenPosition(t *testing.T) {
	tradingRps := new(mocks.TradingServiceRepository)
	tradingRps.On("ClosePosition", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s := NewTradingServiceService(nil, nil, tradingRps)
	s.addInfoToManager(&testPosition)
	go s.OpenPosition(&testPosition)
	s.profMngr.Shares["Tesla"] = 199
	s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].Close <- true

	time.Sleep(time.Millisecond)
	_, ok := s.profMngr.Profiles[testPosition.ProfileID]
	require.True(t, !ok)

	s.addInfoToManager(&testPosition)
	go s.OpenPosition(&testPosition)
	s.profMngr.Shares["Tesla"] = 200

	time.Sleep(time.Millisecond)
	_, ok = s.profMngr.Profiles[testPosition.ProfileID]
	require.True(t, !ok)
}

func TestClosePosition(t *testing.T) {
	tradingRps := new(mocks.TradingServiceRepository)
	tradingRps.On("ClosePosition", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s := NewTradingServiceService(nil, nil, tradingRps)
	s.addInfoToManager(&testPosition)

	s.profMngr.Shares["Tesla"] = 200
	pnl, err := s.ClosePosition(context.Background(), testPosition.ProfileID, testPosition.PositionID)
	require.NoError(t, err)
	require.Equal(t, float64(300), pnl)
}

func TestAddInfoToManager(t *testing.T) {
	s := NewTradingServiceService(nil, nil, nil)
	s.addInfoToManager(&testPosition)
	require.Equal(t, testPosition.PositionID, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].PositionID)
	require.Equal(t, testPosition.ProfileID, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].ProfileID)
	require.Equal(t, testPosition.MoneyAmount, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].MoneyAmount)
	require.Equal(t, testPosition.ShareName, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].ShareName)
	require.Equal(t, testPosition.ShortOrLong, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].ShortOrLong)
	require.Equal(t, testPosition.StopLoss, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].StopLoss)
	require.Equal(t, testPosition.TakeProfit, s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].TakeProfit)
}

func TestCountPnl(t *testing.T) {
	s := NewTradingServiceService(nil, nil, nil)

	pnl := s.countPnl(testPosition.ShareStartPrice, testPosition.ShareEndPrice, testPosition.ShareAmount, testPosition.ShortOrLong)

	require.Equal(t, float64(300), pnl)
}

func TestFillTheCloseChan(t *testing.T) {
	s := NewTradingServiceService(nil, nil, nil)

	s.addInfoToManager(&testPosition)

	s.fillTheCloseChan(testPosition.ProfileID, testPosition.PositionID)
	require.Equal(t, 1, len(s.profMngr.Profiles[testPosition.ProfileID][testPosition.PositionID].Close))
}

func TestCheckStopLossAndTakeProfit(t *testing.T) {
	tradingRps := new(mocks.TradingServiceRepository)
	tradingRps.On("ClosePosition", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s := NewTradingServiceService(nil, nil, tradingRps)
	err := make(chan error, 1)

	s.addInfoToManager(&testPosition)
	s.profMngr.Shares["Apple"] = 200
	s.checkStopLossAndTakeProfit(context.Background(), &testPosition.Position, err)
	require.NoError(t, <-err)
	close(err)

	err = make(chan error, 1)
	testPosition.ShortOrLong = short
	s.checkStopLossAndTakeProfit(context.Background(), &testPosition.Position, err)
	require.NoError(t, <-err)
	close(err)
}

func TestGarbageCleanAfterPosition(t *testing.T) {
	s := NewTradingServiceService(nil, nil, nil)
	s.addInfoToManager(&testPosition)
	testFirstPositionID := testPosition.PositionID
	testPosition.PositionID = uuid.New()
	s.addInfoToManager(&testPosition)
	s.garbageCleanAfterPosition(testPosition.ProfileID, testFirstPositionID)

	_, ok := s.profMngr.Profiles[testPosition.ProfileID][testFirstPositionID]
	require.True(t, !ok)
	_, ok = s.profMngr.Profiles[testPosition.ProfileID]
	require.True(t, ok)

	s.garbageCleanAfterPosition(testPosition.ProfileID, testPosition.PositionID)
	_, ok = s.profMngr.Profiles[testPosition.ProfileID]
	require.True(t, !ok)
}

func TestWaitingForNotification(t *testing.T) {
	profileID, err := uuid.Parse("db85739e-abd7-4a65-9128-7e1528bf0c5a")
	require.NoError(t, err)
	positionID, err := uuid.Parse("971fccd8-d211-4ebb-88c7-0049541d9845")
	require.NoError(t, err)
	balanceRps := new(mocks.BalanceRepository)
	balanceRps.On("WithdrawOnPosition", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	balanceRps.On("MoneyBackWithPnl", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tradingRps := new(mocks.TradingServiceRepository)
	tradingRps.On("ReadPositionRow", mock.Anything, mock.Anything).Return(float64(10), float64(30), float64(1), long, nil)

	s := NewTradingServiceService(nil, balanceRps, tradingRps)
	s.profMngr.Shares["Tesla"] = 200
	l := new(pq.Listener)
	l.Notify = make(chan *pq.Notification)

	go s.WaitingForNotification(context.Background(), l)
	l.Notify <- &pq.Notification{Extra: insertActionPayload}

	time.Sleep(time.Millisecond)
	s.profMngr.Mu.RLock()
	_, ok := s.profMngr.Profiles[profileID][positionID]
	s.profMngr.Mu.RUnlock()

	require.True(t, ok)

	l.Notify <- &pq.Notification{Extra: updateActionPayload}
	time.Sleep(time.Millisecond)
	s.profMngr.Mu.RLock()
	_, ok = s.profMngr.Profiles[profileID]
	s.profMngr.Mu.RUnlock()

	require.True(t, !ok)
}

func TestBackupAllOpenedPositions(t *testing.T) {
	tradingRps := new(mocks.TradingServiceRepository)
	openedPositions := make([]*model.OpenedPosition, 0)
	openedPositions = append(openedPositions, &testPosition)
	tradingRps.On("BackupAllOpenedPositions", mock.Anything).Return(openedPositions, nil)
	s := NewTradingServiceService(nil, nil, tradingRps)
	s.profMngr.Shares["Tesla"] = 199
	s.BackupAllOpenedPositions(context.Background())

	require.Equal(t, 1, len(s.profMngr.Profiles[testPosition.ProfileID]))
}

func TestAddPosition(t *testing.T) {
	tradingRps := new(mocks.TradingServiceRepository)
	tradingRps.On("AddPosition", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s := NewTradingServiceService(nil, nil, tradingRps)
	s.profMngr.Shares[testPosition.ShareName] = 150
	err := s.AddPosition(context.Background(), &testPosition.Position)
	require.NoError(t, err)
}
