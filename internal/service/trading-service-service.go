// Package service contains the bisnes logic of app
package service

import (
	"context"
	"fmt"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
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
	BackupAllOpenedPositions(ctx context.Context) ([]*model.Position, error)
	ReadInfoAboutOpenedPosition(ctx context.Context, positionID uuid.UUID) (float64, float64, error)
}

// TradingServiceService is the struct that implements PriceServiceRepository and TradingServiceRepository interfaces,
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
			ProfilesShares:      make(map[uuid.UUID]chan model.Share),
			Profiles:            make(map[uuid.UUID][]string),
			SharesPositionsRead: make(map[uuid.UUID]map[string]int),
			PositionClose:       make(map[uuid.UUID]chan bool),
		},
	}
}

// SendSharesToProfiles sends shares given from repository PriceServiceRepository to all profiles with opened positions
func (s *TradingServiceService) SendSharesToProfiles(ctx context.Context, length int) {
	var isAnyProfiles bool
	allShares := make(chan model.Share, length)
	errSubscribe := make(chan error, 1)
	shares := make(map[string]float64)
	share := model.Share{}
	for {
		if len(s.profMngr.Profiles) == 0 {
			continue
		}
		if !isAnyProfiles {
			go s.priceRps.Subscribe(ctx, uuid.New(), allShares, errSubscribe)
			isAnyProfiles = true
		}

		err := <-errSubscribe
		if err != nil {
			logrus.WithField("length", length).Errorf("TradingServiceService -> SendSharesToProfiles -> %v", err)
			return
		}

		for i := 0; i < cap(allShares); i++ {
			share = <-allShares
			shares[share.Name] = share.Price
		}
		s.profMngr.Mu.Lock()
		for profileID, selcetedShares := range s.profMngr.Profiles {
			if len(s.profMngr.ProfilesShares[profileID]) != 0 {
				continue
			}
			for _, selectedShare := range selcetedShares {
				select {
				case <-ctx.Done():
					s.profMngr.Mu.Unlock()
					return
				case s.profMngr.ProfilesShares[profileID] <- model.Share{Name: selectedShare, Price: shares[selectedShare]}:
					fmt.Println("INPUT: ", selectedShare, " ", shares[selectedShare])
				}
			}
		}
		s.profMngr.Mu.Unlock()
	}
}

// AddInfoToManager adds new profile to profMngr and if profile already exists just adds new selected share and increments SharesPositionsRead by profileID
func (s *TradingServiceService) AddInfoToManager(ctx context.Context, profileID uuid.UUID, selectedShare string, amount float64, backup bool) error {
	if !backup {
		err := s.balanceRps.WithdrawOnPosition(ctx, profileID.String(), amount)
		if err != nil {
			return fmt.Errorf("TradingServiceService -> AddInfoToManager -> %w", err)
		}
	}

	selectedShares := make([]string, 0)

	if _, ok := s.profMngr.Profiles[profileID]; !ok {
		selectedShares = append(selectedShares, selectedShare)
		s.profMngr.Mu.Lock()
		s.profMngr.Profiles[profileID] = selectedShares
		s.profMngr.ProfilesShares[profileID] = make(chan model.Share, 1)
		s.profMngr.Mu.Unlock()
	} else {
		s.profMngr.Mu.RLock()
		selectedShares = s.profMngr.Profiles[profileID]
		s.profMngr.Mu.RUnlock()
		isSuchShareExist := false
		for _, share := range selectedShares {
			if share == selectedShare {
				isSuchShareExist = true
			}
		}
		if !isSuchShareExist {
			selectedShares = append(selectedShares, selectedShare)
			s.profMngr.Mu.Lock()
			s.profMngr.Profiles[profileID] = selectedShares
			s.profMngr.Mu.Unlock()
		}
		s.profMngr.Mu.Lock()
		s.profMngr.ProfilesShares[profileID] = make(chan model.Share, cap(s.profMngr.ProfilesShares[profileID])+1)
		s.profMngr.Mu.Unlock()
	}

	if _, ok := s.profMngr.SharesPositionsRead[profileID]; !ok {
		s.profMngr.Mu.Lock()
		s.profMngr.SharesPositionsRead[profileID] = make(map[string]int)
		s.profMngr.Mu.Unlock()
	}

	s.profMngr.SharesPositionsRead[profileID][selectedShare]++

	return nil
}

// DeleteProfileAtAll deletes profile from Profiles map and ProfilesShares in ProfileManager by profileID
func (s *TradingServiceService) DeleteProfileAtAll(profileID uuid.UUID) {
	s.profMngr.Mu.Lock()
	defer s.profMngr.Mu.Unlock()
	if _, ok := s.profMngr.Profiles[profileID]; ok {
		delete(s.profMngr.Profiles, profileID)
		close(s.profMngr.ProfilesShares[profileID])
		delete(s.profMngr.ProfilesShares, profileID)
		delete(s.profMngr.SharesPositionsRead, profileID)
	}
}

// CountPnl counts pnl of position using share amount, share start and end prices
func (s *TradingServiceService) CountPnl(vector string, shareAmount, shareStartPrice, shareEndPrice float64) float64 {
	shareStartPriceDecimal := decimal.NewFromFloat(shareStartPrice)
	shareEndPriceDecimal := decimal.NewFromFloat(shareEndPrice)
	shareAmountDecimal := decimal.NewFromFloat(shareAmount)
	switch vector {
	case long:
		temp := shareAmountDecimal.Mul(shareEndPriceDecimal.Sub(shareStartPriceDecimal))
		return temp.InexactFloat64()
	case short:
		temp := shareAmountDecimal.Mul(shareStartPriceDecimal.Sub(shareEndPriceDecimal))
		return temp.InexactFloat64()
	}
	return 0
}

// ClosePosition fills the initialized with position ID map PositionClose in profile manager with true value
func (s *TradingServiceService) ClosePosition(positionID uuid.UUID) error {
	s.profMngr.Mu.RLock()
	if _, ok := s.profMngr.PositionClose[positionID]; !ok {
		s.profMngr.Mu.RUnlock()
		return fmt.Errorf("TradingServiceService -> ClosePosition -> error: position with such ID was never opened")
	}
	s.profMngr.Mu.RUnlock()

	s.profMngr.Mu.Lock()
	s.profMngr.PositionClose[positionID] <- true
	s.profMngr.Mu.Unlock()
	return nil
}

// chooseBetweenVector checks if actual share price is in the borders of stoploss and takeprofit of exact position
func (s *TradingServiceService) chooseBetweenVector(position *model.Position, share *model.Share) {
	switch position.Vector {
	case long:
		if share.Name == position.ShareName && (share.Price >= position.TakeProfit || share.Price <= position.StopLoss) {
			s.profMngr.PositionClose[position.PositionID] <- true
		}
	case short:
		if share.Name == position.ShareName && (share.Price >= position.StopLoss || share.Price <= position.TakeProfit) {
			s.profMngr.PositionClose[position.PositionID] <- true
		}
	}
}

// CheckStopLossAndTakeProfit called in goroutine with fill the given channel when position woult be at StopLoss or TakeProfit
func (s *TradingServiceService) CheckStopLossAndTakeProfit(position *model.Position, outOfPosition chan bool) {
	shares := make(map[string]float64)
	share := model.Share{}
	for {
		s.profMngr.Mu.Lock()
		if len(s.profMngr.ProfilesShares[position.ProfileID]) != len(s.profMngr.Profiles[position.ProfileID]) {
			s.profMngr.Mu.Unlock()
			continue
		}
		for i := 0; i < cap(s.profMngr.ProfilesShares[position.ProfileID]); i++ {
			share = <-s.profMngr.ProfilesShares[position.ProfileID]
			shares[share.Name] = share.Price
		}
		s.profMngr.Mu.Unlock()
		for key, value := range shares {
			if position.ShareName != key {
				continue
			}
			select {
			case <-outOfPosition:
				return
			default:
				s.chooseBetweenVector(position, &model.Share{Name: key, Price: value})
				if len(s.profMngr.PositionClose[position.PositionID]) > 0 {
					return
				}
				fmt.Println("CHECK: ", share.Name, " ", share.Price)
			}
		}
	}
}

// garbageCleanAfterPosition delete from profile manager info that was used by exact position
func (s *TradingServiceService) garbageCleanAfterPosition(position *model.Position) {
	s.profMngr.Mu.Lock()
	s.profMngr.SharesPositionsRead[position.ProfileID][position.ShareName]--
	if s.profMngr.SharesPositionsRead[position.ProfileID][position.ShareName] == 0 {
		delete(s.profMngr.SharesPositionsRead[position.ProfileID], position.ShareName)
		for i, share := range s.profMngr.Profiles[position.ProfileID] {
			if share == position.ShareName {
				s.profMngr.Profiles[position.ProfileID] = append(s.profMngr.Profiles[position.ProfileID][:i], s.profMngr.Profiles[position.ProfileID][i+1:]...)
				break
			}
		}
	}
	s.profMngr.Mu.Unlock()
	if len(s.profMngr.SharesPositionsRead[position.ProfileID]) == 0 {
		s.DeleteProfileAtAll(position.ProfileID)
	} else {
		s.profMngr.Mu.Lock()
		close(s.profMngr.ProfilesShares[position.ProfileID])
		delete(s.profMngr.ProfilesShares, position.PositionID)
		s.profMngr.ProfilesShares[position.ProfileID] = make(chan model.Share, cap(s.profMngr.ProfilesShares[position.ProfileID])-1)
		s.profMngr.Mu.Unlock()
	}

	s.profMngr.Mu.Lock()
	close(s.profMngr.PositionClose[position.PositionID])
	delete(s.profMngr.PositionClose, position.PositionID)
	s.profMngr.Mu.Unlock()
}

// BackupAllOpenedPositions adds info to manager about all positions that wasn't be closed and calls the method OpenPosition to each one
func (s *TradingServiceService) BackupAllOpenedPositions(ctx context.Context) {
	openedPositions, err := s.tradingRps.BackupAllOpenedPositions(ctx)
	if err != nil {
		logrus.Errorf("TradingServiceService -> BackupAllOpenedPositions -> %v", err)
	}
	for _, position := range openedPositions {
		err := s.AddInfoToManager(ctx, position.ProfileID, position.ShareName, position.MoneyAmount, true)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"position.ProfileID":   position.ProfileID,
				"position.ShareName":   position.ShareName,
				"position.MoneyAmount": position.MoneyAmount,
			}).Errorf("TradingServiceService -> BackupAllOpenedPositions -> %v", err)
		}
		go func(position *model.Position) {
			_, err := s.OpenPosition(ctx, position, true)
			logrus.WithFields(logrus.Fields{
				"position": position,
				"backup":   true,
			}).Errorf("TradingServiceService -> BackupAllOpenedPositions -> %v", err)
		}(position)
	}
}

// takeActualShares takes actual shares and theirs prices for an exact profileID and return it in map of share names as a key and prices as value
func (s *TradingServiceService) takeActualShares(profileID uuid.UUID) map[string]float64 {
	shares := make(map[string]float64)

	for {
		s.profMngr.Mu.Lock()

		if len(s.profMngr.ProfilesShares[profileID]) != len(s.profMngr.Profiles[profileID]) {
			s.profMngr.Mu.Unlock()
			continue
		}

		share := <-s.profMngr.ProfilesShares[profileID]

		shares[share.Name] = share.Price
		for i := 1; i < cap(s.profMngr.ProfilesShares[profileID]); i++ {
			share = <-s.profMngr.ProfilesShares[profileID]
			shares[share.Name] = share.Price
		}

		s.profMngr.Mu.Unlock()

		break
	}

	return shares
}

// OpenPosition calls method AddPosition from TradingServiceRepository and wait till the end of position to return the its pnl
func (s *TradingServiceService) OpenPosition(ctx context.Context, position *model.Position, backup bool) (float64, error) {
	moneyAmountDecimal := decimal.NewFromFloat(position.MoneyAmount)
	var pnl float64
	var shareStartPrice float64
	var shareAmount float64
	var moneyBack decimal.Decimal
	var err error

	if !backup {
		shares := s.takeActualShares(position.ProfileID)

		shareStartPrice = shares[position.ShareName]
		shareStartPriceDecimal := decimal.NewFromFloat(shareStartPrice)
		shareAmount = moneyAmountDecimal.Div(shareStartPriceDecimal).InexactFloat64()

		err = s.tradingRps.AddPosition(ctx, position, shareAmount, shareStartPrice)
		if err != nil {
			return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
		}
	} else {
		shareStartPrice, shareAmount, err = s.tradingRps.ReadInfoAboutOpenedPosition(ctx, position.PositionID)
		if err != nil {
			return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
		}
	}

	outOfPosition := make(chan bool, 1)
	s.profMngr.PositionClose[position.PositionID] = make(chan bool, 1)
	go s.CheckStopLossAndTakeProfit(position, outOfPosition)
	shares := make(map[string]float64)
	select {
	case <-s.profMngr.PositionClose[position.PositionID]:
		outOfPosition <- true
		for {
			s.profMngr.Mu.Lock()
			if len(s.profMngr.ProfilesShares[position.ProfileID]) != len(s.profMngr.Profiles[position.ProfileID]) {
				s.profMngr.Mu.Unlock()
				continue
			}
			for i := 0; i < cap(s.profMngr.ProfilesShares[position.ProfileID]); i++ {
				share := <-s.profMngr.ProfilesShares[position.ProfileID]
				shares[share.Name] = share.Price
			}
			s.profMngr.Mu.Unlock()
			break
		}

	case <-outOfPosition:
		s.profMngr.Mu.Lock()
		for i := 0; i < cap(s.profMngr.ProfilesShares[position.ProfileID]); i++ {
			share := <-s.profMngr.ProfilesShares[position.ProfileID]
			shares[share.Name] = share.Price
		}
		s.profMngr.Mu.Unlock()
	}

	if _, ok := shares[position.ShareName]; !ok {
		return 0, fmt.Errorf("TradingServiceService -> ClosePosition -> error: position on such share doesn't exist")
	}

	shareEndPrice := shares[position.ShareName]
	err = s.tradingRps.ClosePosition(ctx, position.PositionID, shareEndPrice)
	if err != nil {
		return 0, fmt.Errorf("TradingServiceService -> ClosePosition -> %w", err)
	}

	pnl = s.CountPnl(position.Vector, shareAmount, shareStartPrice, shareEndPrice)

	pnlDecimal := decimal.NewFromFloat(pnl)
	moneyBack = pnlDecimal.Add(moneyAmountDecimal)

	s.garbageCleanAfterPosition(position)

	err = s.balanceRps.MoneyBackWithPnl(ctx, position.ProfileID.String(), moneyBack.InexactFloat64())
	if err != nil {
		return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
	}

	return pnl, nil
}

// ReadAllOpenedPositionsByProfileID calls similar method of trading service repository and returnes slice higher to handler
func (s *TradingServiceService) ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error) {
	return s.tradingRps.ReadAllOpenedPositionsByProfileID(ctx, profileID)
}
