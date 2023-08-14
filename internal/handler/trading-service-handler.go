// Package handler contains methods that handle requests and send responses
package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/distuurbia/tradingService/internal/config"
	"github.com/distuurbia/tradingService/internal/model"
	protocol "github.com/distuurbia/tradingService/protocol/trading"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TradingServiceService is an interface of service TradingServiceService
type TradingServiceService interface {
	AddInfoToManager(ctx context.Context, profileID uuid.UUID, selectedShare string, amount float64) error
	OpenPosition(ctx context.Context, position *model.Position) (float64, error)
	ClosePosition(ctx context.Context, profileID, positionID uuid.UUID) error
}

// TradingServiceHandler is an object that implement TradingServiceService and also contains an object of type *validator.Validate
type TradingServiceHandler struct {
	s        TradingServiceService
	validate *validator.Validate
	cfg      *config.Config
	protocol.UnimplementedTradingServiceServiceServer
}

// NewTradingServiceHandler is a constructor for TradingServiceHandler
func NewTradingServiceHandler(s TradingServiceService, validate *validator.Validate, cfg *config.Config) *TradingServiceHandler {
	return &TradingServiceHandler{s: s, validate: validate, cfg: cfg}
}

// ValidationID validate given in and parses it to uuid.UUID type
func (h *TradingServiceHandler) ValidationID(ctx context.Context, id string) (uuid.UUID, error) {
	err := h.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("ValidationID -> %v", err)
		return uuid.Nil, err
	}

	profileID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("ValidationID -> %v", err)
		return uuid.Nil, err
	}

	if profileID == uuid.Nil {
		logrus.Errorf("ValidationID -> error: failed to use uuid")
		return uuid.Nil, fmt.Errorf("ValidationID -> error: failed to use uuid")
	}
	return profileID, nil
}

// OpenPosition opens position by calling method service TradingServiceService OpenPosition
func (h *TradingServiceHandler) OpenPosition(ctx context.Context, req *protocol.OpenPositionRequest) (*protocol.OpenPositionResponse, error) {
	profileID, err := h.ValidationID(ctx, req.Position.ProfileID)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}
	positionID, err := h.ValidationID(ctx, req.Position.PositionID)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}
	position := model.Position{
		ProfileID:   profileID,
		PositionID:  positionID,
		StopLoss:    req.Position.StopLoss,
		TakeProfit:  req.Position.TakeProfit,
		MoneyAmount: req.Position.MoneyAmount,
		ShareName:   req.Position.ShareName,
		Vector:      req.Position.Vector,
	}
	err = h.validate.StructCtx(ctx, position)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}

	if !strings.Contains(h.cfg.TradingServiceShares, position.ShareName) {
		logrus.Errorf("TradingServiceHandler -> OpenPosition -> error: such share doesn't exist")
		return &protocol.OpenPositionResponse{}, fmt.Errorf("TradingServiceHandler -> OpenPosition -> error: such share doesn't exist")
	}

	err = h.s.AddInfoToManager(ctx, position.ProfileID, position.ShareName, position.MoneyAmount)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}

	pnl, err := h.s.OpenPosition(ctx, &position)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}
	return &protocol.OpenPositionResponse{Pnl: pnl}, nil
}

// ClosePosition closes position by calling method service TradingServiceService ClosePosition
func (h *TradingServiceHandler) ClosePosition(ctx context.Context, req *protocol.ClosePositionRequest) (*protocol.ClosePositionResponse, error) {
	profileID, err := h.ValidationID(ctx, req.ProfileID)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}
	positionID, err := h.ValidationID(ctx, req.PositionID)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}

	err = h.s.ClosePosition(ctx, profileID, positionID)
	if err != nil {
		logrus.Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}
	return &protocol.ClosePositionResponse{}, nil
}

// CheckTradableShares returns what shares trading service subscribing
func (h *TradingServiceHandler) CheckTradableShares(_ context.Context, _ *protocol.CheckTradableSharesRequest) (*protocol.CheckTradableSharesResponse, error) {
	return &protocol.CheckTradableSharesResponse{TradableShares: h.cfg.TradingServiceShares}, nil
}
