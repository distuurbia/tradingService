// Package service contains the bisnes logic of app
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// long is the long vector and short is the short vector
const (
	long  = "long"
	short = "short"
)

// PriceServiceRepository is an interface of repository PriceServiceRepository
type PriceServiceRepository interface {
	Subscribe(ctx context.Context, subID uuid.UUID, subscribersShares chan model.Share, errSubscribe chan error)
}

// BalanceRepository is an interface of repository BalanceRepository
type BalanceRepository interface {
	WithdrawOnPosition(ctx context.Context, profileID string, amount float64) error
	MoneyBackWithPnl(ctx context.Context, profileID string, pnl float64) error
}

// TradingServiceRepository is an interface of repository TradingServiceRepository
type TradingServiceRepository interface {
	AddPosition(ctx context.Context, position *model.Position, shareAmount, shareStartPrice float64) error
	ClosePosition(ctx context.Context, positionID uuid.UUID, shareEndPrice float64) error
	ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error)
	BackupAllOpenedPositions(ctx context.Context) ([]*model.OpenedPosition, error)
	ReadPositionRow(ctx context.Context, positionID uuid.UUID) (shareStartPrice, shareEndPrice,
		shareAmount float64, shortOrLong string, err error)
	ReadShareNameByPositionID(ctx context.Context, positionID uuid.UUID) (string, error)
}

// TradingServiceService implements PriceServiceRepository, BalanceRepository and TradingServiceRepository interfaces,
// also contains the object of *model.ProfilesManager
type TradingServiceService struct {
	priceRps   PriceServiceRepository
	balanceRps BalanceRepository
	tradingRps TradingServiceRepository
	profMngr   *model.ProfilesManager
}

// NewTradingServiceService is the constructor of TradingServiceService struct
func NewTradingServiceService(priceRps PriceServiceRepository, balanceRps BalanceRepository, tradingRps TradingServiceRepository) *TradingServiceService {
	return &TradingServiceService{
		priceRps:   priceRps,
		balanceRps: balanceRps,
		tradingRps: tradingRps,
		profMngr: &model.ProfilesManager{
			Shares:   make(map[string]float64),
			Profiles: make(map[uuid.UUID]map[uuid.UUID]model.OpenedPosition),
		},
	}
}

// SendSharesToProfiles updates shares given from repository PriceServiceRepository in map Shares of profMngr
func (s *TradingServiceService) SendSharesToProfiles(ctx context.Context, length int) {
	allShares := make(chan model.Share, length)
	errSubscribe := make(chan error, 1)
	go s.priceRps.Subscribe(ctx, uuid.New(), allShares, errSubscribe)
	for {
		err := <-errSubscribe
		if err != nil {
			logrus.WithField("length", length).Errorf("TradingServiceService -> SendSharesToProfiles -> %v", err)
			return
		}

		s.profMngr.Mu.Lock()
		for i := 0; i < cap(allShares); i++ {
			share := <-allShares
			s.profMngr.Shares[share.Name] = share.Price
		}
		s.profMngr.Mu.Unlock()
	}
}

// AddInfoToManager adds new position for an exact profile to profMngr
func (s *TradingServiceService) addInfoToManager(position *model.OpenedPosition) {
	s.profMngr.Mu.Lock()
	if _, ok := s.profMngr.Profiles[position.ProfileID]; !ok {
		s.profMngr.Profiles[position.ProfileID] = make(map[uuid.UUID]model.OpenedPosition)
	}
	position.Close = make(chan bool, 1)
	s.profMngr.Profiles[position.ProfileID][position.PositionID] = *position
	s.profMngr.Mu.Unlock()
}

// CountPnl counts pnl of position using its positionID
func (s *TradingServiceService) countPnl(shareStartPrice, shareEndPrice, shareAmount float64, shortOrLong string) float64 {
	shareStartPriceDecimal := decimal.NewFromFloat(shareStartPrice)
	shareEndPriceDecimal := decimal.NewFromFloat(shareEndPrice)
	shareAmountDecimal := decimal.NewFromFloat(shareAmount)
	pnlDecimal := decimal.Decimal{}
	switch shortOrLong {
	case long:
		pnlDecimal = shareAmountDecimal.Mul(shareEndPriceDecimal.Sub(shareStartPriceDecimal))
	case short:
		pnlDecimal = shareAmountDecimal.Mul(shareStartPriceDecimal.Sub(shareEndPriceDecimal))
	}

	return pnlDecimal.InexactFloat64()
}

// FillTheCloseChan fills the Close chan of an exact position by its ID
func (s *TradingServiceService) fillTheCloseChan(profileID, positionID uuid.UUID) {
	s.profMngr.Mu.Lock()
	s.profMngr.Profiles[profileID][positionID].Close <- true
	s.profMngr.Mu.Unlock()
}

// CheckStopLossAndTakeProfit checks if actual price of share of an exact position is in the boards of StopLoss and TakeProfit
func (s *TradingServiceService) checkStopLossAndTakeProfit(ctx context.Context, position *model.Position, err chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		s.profMngr.Mu.RLock()
		if _, ok := s.profMngr.Profiles[position.ProfileID][position.PositionID]; !ok {
			s.profMngr.Mu.RUnlock()
			continue
		}
		sharePrice := s.profMngr.Shares[position.ShareName]
		s.profMngr.Mu.RUnlock()
		switch position.ShortOrLong {
		case long:
			if sharePrice >= position.TakeProfit || sharePrice <= position.StopLoss {
				errClose := s.tradingRps.ClosePosition(ctx, position.PositionID, sharePrice)
				if errClose != nil {
					err <- fmt.Errorf("TradingServiceService -> CheckStopLossAndTakeProfit -> %w", errClose)
					return
				}
				err <- nil
				return
			}
		case short:
			if sharePrice >= position.StopLoss || sharePrice <= position.TakeProfit {
				errClose := s.tradingRps.ClosePosition(ctx, position.PositionID, sharePrice)
				if errClose != nil {
					err <- fmt.Errorf("TradingServiceService -> CheckStopLossAndTakeProfit -> %w", errClose)
					return
				}
				err <- nil
				return
			}
		}
	}
}

// garbageCleanAfterPosition closes Close chan of an exact Position and deletes this position from Profiles map in profMngr
func (s *TradingServiceService) garbageCleanAfterPosition(profileID, positionID uuid.UUID) {
	s.profMngr.Mu.Lock()
	close(s.profMngr.Profiles[profileID][positionID].Close)
	delete(s.profMngr.Profiles[profileID], positionID)

	if len(s.profMngr.Profiles[profileID]) == 0 {
		delete(s.profMngr.Profiles, profileID)
	}
	s.profMngr.Mu.Unlock()
}

// WaitingForNotification waits for notification from channel l.Notify and after exact actions in db does some logic
func (s *TradingServiceService) WaitingForNotification(ctx context.Context, l *pq.Listener) {
	var valueToParse = struct {
		Action               string `json:"action"`
		model.OpenedPosition `json:"data"`
	}{}
	for {
		notification := <-l.Notify

		err := json.Unmarshal([]byte(notification.Extra), &valueToParse)
		logrus.Info(notification)
		if err != nil {
			logrus.WithField("payload", notification.Extra).Errorf("TradingServiceService -> WaitingForNotification -> %v", err)
			continue
		}
		if valueToParse.Action == "INSERT" {
			shareAmountDecimal := decimal.NewFromFloat(valueToParse.ShareAmount)
			shareStartPriceDecimal := decimal.NewFromFloat(valueToParse.ShareStartPrice)
			valueToParse.MoneyAmount = shareAmountDecimal.Mul(shareStartPriceDecimal).InexactFloat64()
			errInsert := s.balanceRps.WithdrawOnPosition(ctx, valueToParse.ProfileID.String(), valueToParse.MoneyAmount)
			if errInsert != nil {
				logrus.WithField("position.PositionID", valueToParse.PositionID).Errorf("TradingServiceService -> WaitingForNotification -> %v", err)
			}
			s.addInfoToManager(&valueToParse.OpenedPosition)
			go s.OpenPosition(&valueToParse.OpenedPosition)
		} else if valueToParse.Action == "UPDATE" {
			pnl := s.countPnl(valueToParse.ShareStartPrice, valueToParse.ShareEndPrice, valueToParse.ShareAmount, valueToParse.ShortOrLong)

			pnlDecimal := decimal.NewFromFloat(pnl)
			moneyAmountDecimal := decimal.NewFromFloat(valueToParse.MoneyAmount)
			moneyBack := pnlDecimal.Add(moneyAmountDecimal)
			err = s.balanceRps.MoneyBackWithPnl(ctx, valueToParse.ProfileID.String(), moneyBack.InexactFloat64())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"valueToParse.ProfileID": valueToParse.ProfileID,
					"moneyBack":              moneyBack.InexactFloat64(),
				}).Errorf("TradingServiceService -> WaitingForNotification -> %v", err)
			}
			s.fillTheCloseChan(valueToParse.ProfileID, valueToParse.PositionID)
		}
	}
}

// BackupAllOpenedPositions adds info to manager about all positions that wasn't be closed and calls methods AddInfoToManager and OpenPosition to each one
func (s *TradingServiceService) BackupAllOpenedPositions(ctx context.Context) {
	openedPositions, err := s.tradingRps.BackupAllOpenedPositions(ctx)
	if err != nil {
		logrus.Errorf("TradingServiceService -> BackupAllOpenedPositions -> %v", err)
	}
	for _, position := range openedPositions {
		s.addInfoToManager(position)
		go s.OpenPosition(position)
	}
}

// OpenPosition calls method CheckStopLossAndTakeProfit in goroutine and waits until chan Close of an exact position or an err won't be fullfiled
// after that calls method garbageCleanAfterPosition to delete info about already closed position from manager
func (s *TradingServiceService) OpenPosition(position *model.OpenedPosition) {
	errChan := make(chan error, 1)
	ctxCheck, cancel := context.WithCancel(context.Background())
	go s.checkStopLossAndTakeProfit(ctxCheck, &position.Position, errChan)
	select {
	case <-s.profMngr.Profiles[position.ProfileID][position.PositionID].Close:
		cancel()
	case err := <-errChan:
		cancel()
		if err != nil {
			logrus.WithField("position", position).Errorf("TradingServiceService -> OpenPosition -> %v", err)
		}
	}

	s.garbageCleanAfterPosition(position.ProfileID, position.PositionID)
}

// ReadAllOpenedPositionsByProfileID calls similar method of trading service repository and returnes slice higher to  the handler
func (s *TradingServiceService) ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error) {
	return s.tradingRps.ReadAllOpenedPositionsByProfileID(ctx, profileID)
}

// AddPosition adds info about position to db
func (s *TradingServiceService) AddPosition(ctx context.Context, position *model.Position) error {
	s.profMngr.Mu.RLock()
	shareStartPrice := s.profMngr.Shares[position.ShareName]
	s.profMngr.Mu.RUnlock()

	moneyAmountDecimal := decimal.NewFromFloat(position.MoneyAmount)
	shareStartPriceDecimal := decimal.NewFromFloat(shareStartPrice)
	shareAmount := moneyAmountDecimal.Div(shareStartPriceDecimal).InexactFloat64()

	err := s.tradingRps.AddPosition(ctx, position, shareAmount, shareStartPrice)
	if err != nil {
		return fmt.Errorf("TradingServiceService -> AddPosition -> %w", err)
	}

	return nil
}

// ClosePosition adds to db row for an exact position closeTime and shareEndPrice values also returns pnl value
func (s *TradingServiceService) ClosePosition(ctx context.Context, profileID, positionID uuid.UUID) (float64, error) {
	s.profMngr.Mu.RLock()
	shareName := s.profMngr.Profiles[profileID][positionID].ShareName
	shareEndPrice := s.profMngr.Shares[shareName]
	s.profMngr.Mu.RUnlock()
	err := s.tradingRps.ClosePosition(ctx, positionID, shareEndPrice)
	if err != nil {
		return 0, fmt.Errorf("TradingServiceService -> ClosePosition -> %w", err)
	}
	s.profMngr.Mu.RLock()
	shareStartPrice := s.profMngr.Profiles[profileID][positionID].ShareStartPrice
	shareAmount := s.profMngr.Profiles[profileID][positionID].ShareAmount
	shortOrLong := s.profMngr.Profiles[profileID][positionID].ShortOrLong
	s.profMngr.Mu.RUnlock()

	pnl := s.countPnl(shareStartPrice, shareEndPrice, shareAmount, shortOrLong)

	return pnl, nil
}
