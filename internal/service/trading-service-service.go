// Package service contains the bisnes logic of app
package service

import (
	"context"
	"fmt"
	"strings"

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
	Subscribe(ctx context.Context, subID uuid.UUID, subscribersShares chan []*model.Share, errSubscribe chan error)
}

// BalanceRepository is an interface of repository BalanceRepository
type BalanceRepository interface {
	WithdrawOnPosition(ctx context.Context, profileID string, amount float64) error
	MoneyBackWithPnl(ctx context.Context, profileID string, pnl float64) error
}

// TradingServiceRepository is an interface of repository TradingServiceRepository
type TradingServiceRepository interface {
	AddPosition(ctx context.Context, position *model.Position, shareStartPrice float64) error
	ClosePosition(ctx context.Context, positionID uuid.UUID, shareEndPrice float64) error
	ReadShareNameByID(ctx context.Context, positionID uuid.UUID) (string, error)
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
			ProfilesShares:      make(map[uuid.UUID]chan []*model.Share),
			Profiles:            make(map[uuid.UUID][]string),
			SharesPositionsRead: make(map[uuid.UUID]map[string]int),
			PositionClose:       make(map[uuid.UUID]chan bool),
		},
	}
}

// SendSharesToProfiles sends shares given from repository PriceServiceRepository to all profiles with opened positions
func (s *TradingServiceService) SendSharesToProfiles(ctx context.Context) {
	var isAnyProfiles bool
	allShares := make(chan []*model.Share)
	errSubscribe := make(chan error)
	for {
		if len(s.profMngr.Profiles) == 0 {
			continue
		}
		if !isAnyProfiles {
			go s.priceRps.Subscribe(ctx, uuid.New(), allShares, errSubscribe)
			isAnyProfiles = true
		}
		if len(s.profMngr.Profiles) == 0 && isAnyProfiles {
			isAnyProfiles = false
			ctx.Done()
			continue
		}
		err := <-errSubscribe
		if err != nil {
			logrus.Infof("TradingServiceService -> SendSharesToProfiles -> %v", err)
		}
		shares := <-allShares
		for subID, selcetedShares := range s.profMngr.Profiles {
			tempShares := make([]*model.Share, 0)

			for _, share := range shares {
				if strings.Contains(strings.Join(selcetedShares, ","), share.Name) {
					tempShares = append(tempShares, share)
				}
			}

			select {
			case <-ctx.Done():
				return
			case s.profMngr.ProfilesShares[subID] <- tempShares:
				fmt.Println("INPUT\n", subID, "\n", tempShares[0].Name, "\n", tempShares[0].Price, "INPUT")
			}
		}
	}
}

// AddInfoToManager adds new profile to profMngr and if profile already exists just adds new selected share and increments SharesPositionsRead by profileID
func (s *TradingServiceService) AddInfoToManager(ctx context.Context, profileID uuid.UUID, selectedShare string, amount float64) error {
	err := s.balanceRps.WithdrawOnPosition(ctx, profileID.String(), amount)
	if err != nil {
		return fmt.Errorf("TradingServiceService -> AddInfoToManager -> %w", err)
	}
	selectedShares := make([]string, 0)
	const msgs = 1

	if _, ok := s.profMngr.Profiles[profileID]; !ok {
		selectedShares = append(selectedShares, selectedShare)
		s.profMngr.Profiles[profileID] = selectedShares
		s.profMngr.ProfilesShares[profileID] = make(chan []*model.Share, msgs)
	} else {
		selectedShares := s.profMngr.Profiles[profileID]
		isSuchShareExist := false
		for _, share := range selectedShares {
			if share == selectedShare {
				isSuchShareExist = true
			}
		}
		if !isSuchShareExist {
			selectedShares = append(selectedShares, selectedShare)
			s.profMngr.Profiles[profileID] = selectedShares
		}
	}
	if _, ok := s.profMngr.SharesPositionsRead[profileID]; !ok {
		s.profMngr.SharesPositionsRead[profileID] = make(map[string]int)
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

// ClosePosition closes exact position and deletes everything related from all included structures in ProfileManager,
// if profile don't have any opened positions, calls DeleteProfileAtAll method
func (s *TradingServiceService) ClosePosition(ctx context.Context, profileID, positionID uuid.UUID) error {
	_, ok := s.profMngr.Profiles[profileID]
	_, secOk := s.profMngr.SharesPositionsRead[profileID]
	if !ok && secOk {
		return fmt.Errorf("TradingServiceService -> ClosePosition -> error: position with such ID doesn't exist")
	}
	var shareEndPrice float64
	shares := <-s.profMngr.ProfilesShares[profileID]
	selectedShare, err := s.tradingRps.ReadShareNameByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("TradingServiceService -> ClosePosition -> %w", err)
	}
	for _, share := range shares {
		if share.Name == selectedShare {
			shareEndPrice = share.Price
		}
	}
	s.profMngr.SharesPositionsRead[profileID][selectedShare]--
	if s.profMngr.SharesPositionsRead[profileID][selectedShare] == 0 {
		delete(s.profMngr.SharesPositionsRead[profileID], selectedShare)
		for i, share := range s.profMngr.Profiles[profileID] {
			if share == selectedShare {
				s.profMngr.Profiles[profileID] = append(s.profMngr.Profiles[profileID][:i], s.profMngr.Profiles[profileID][:i+1]...)
				break
			}
		}
	}

	err = s.tradingRps.ClosePosition(ctx, positionID, shareEndPrice)
	if err != nil {
		return fmt.Errorf("TradingServiceService -> ClosePosition -> %w", err)
	}
	s.profMngr.PositionClose[positionID] <- true
	return nil
}

// CheckStopLossAndTakeProfit called in goroutine with fill the given channel when position woult be at StopLoss or TakeProfit
func (s *TradingServiceService) CheckStopLossAndTakeProfit(position *model.Position, outOfPosition chan bool) {
	for {
		shares := <-s.profMngr.ProfilesShares[position.ProfileID]
		for _, share := range shares {
			switch position.Vector {
			case long:
				if share.Name == position.ShareName && (share.Price >= position.TakeProfit || share.Price <= position.StopLoss) {
					outOfPosition <- true
					return
				}
			case short:
				if share.Name == position.ShareName && (share.Price >= position.StopLoss || share.Price <= position.TakeProfit) {
					outOfPosition <- true
					return
				}
			}
		}
	}
}

// OpenPosition calls method AddPosition from TradingServiceRepository and wait till the end of position to return the its pnl
func (s *TradingServiceService) OpenPosition(ctx context.Context, position *model.Position) (float64, error) {
	shares := <-s.profMngr.ProfilesShares[position.ProfileID]
	var shareStartPrice float64
	var shareAmount float64
	moneyAmountDecimal := decimal.NewFromFloat(position.MoneyAmount)
	var pnl float64
	for _, share := range shares {
		if share.Name != position.ShareName {
			continue
		}
		shareStartPrice = share.Price
		shareStartPriceDecimal := decimal.NewFromFloat(shareStartPrice)
		shareAmount = moneyAmountDecimal.Div(shareStartPriceDecimal).InexactFloat64()
		err := s.tradingRps.AddPosition(ctx, position, shareStartPrice)
		if err != nil {
			return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
		}
	}

	outOfPosition := make(chan bool)
	s.profMngr.PositionClose[position.PositionID] = make(chan bool, 1)
	go s.CheckStopLossAndTakeProfit(position, outOfPosition)

	select {
	case <-s.profMngr.PositionClose[position.PositionID]:
		close(s.profMngr.PositionClose[position.PositionID])
		delete(s.profMngr.PositionClose, position.PositionID)
		shares := <-s.profMngr.ProfilesShares[position.ProfileID]
		if len(s.profMngr.SharesPositionsRead[position.ProfileID]) == 0 {
			s.DeleteProfileAtAll(position.ProfileID)
		}
		for _, share := range shares {
			if share.Name != position.ShareName {
				continue
			}
			pnl = s.CountPnl(position.Vector, shareAmount, shareStartPrice, share.Price)

			pnlDecimal := decimal.NewFromFloat(pnl)
			moneyBack := pnlDecimal.Add(moneyAmountDecimal)
			err := s.balanceRps.MoneyBackWithPnl(ctx, position.ProfileID.String(), moneyBack.InexactFloat64())
			if err != nil {
				return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
			}
			return pnl, nil
		}
	case <-outOfPosition:
		shares := <-s.profMngr.ProfilesShares[position.ProfileID]
		for _, share := range shares {
			if share.Name == position.ShareName {
				pnl = s.CountPnl(position.Vector, shareAmount, shareStartPrice, share.Price)
			}
		}
		err := s.ClosePosition(ctx, position.ProfileID, position.PositionID)
		if err != nil {
			return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
		}
		if len(s.profMngr.SharesPositionsRead[position.ProfileID]) == 0 {
			s.DeleteProfileAtAll(position.ProfileID)
		}
		close(s.profMngr.PositionClose[position.PositionID])
		delete(s.profMngr.PositionClose, position.PositionID)
		pnlDecimal := decimal.NewFromFloat(pnl)
		moneyBack := pnlDecimal.Add(moneyAmountDecimal)
		err = s.balanceRps.MoneyBackWithPnl(ctx, position.ProfileID.String(), moneyBack.InexactFloat64())
		if err != nil {
			return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> %w", err)
		}
		return pnl, nil
	}
	return 0, fmt.Errorf("TradingServiceService -> OpenPosition -> error: failed to return pnl")
}
