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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TradingServiceService is an interface of service TradingServiceService
type TradingServiceService interface {
	AddInfoToManager(ctx context.Context, profileID uuid.UUID, selectedShare string, amount float64, backup bool) error
	OpenPosition(ctx context.Context, position *model.Position, backup bool) (float64, error)
	ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error)
	ClosePosition(positionID uuid.UUID) error
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

// validationID validate given in and parses it to uuid.UUID type
func (h *TradingServiceHandler) validationID(ctx context.Context, id string) (uuid.UUID, error) {
	err := h.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		return uuid.Nil, err
	}

	profileID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	if profileID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("validationID -> error: failed to use uuid")
	}
	return profileID, nil
}

// OpenPosition opens position by calling method service TradingServiceService OpenPosition
func (h *TradingServiceHandler) OpenPosition(ctx context.Context, req *protocol.OpenPositionRequest) (*protocol.OpenPositionResponse, error) {
	profileID, err := h.validationID(ctx, req.Position.ProfileID)
	if err != nil {
		logrus.WithField("req.Position.ProfileID", req.Position.ProfileID).Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}
	positionID, err := h.validationID(ctx, req.Position.PositionID)
	if err != nil {
		logrus.WithField("req.Position.PositionID", req.Position.PositionID).Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
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
		logrus.WithField("position", position).Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}

	if !strings.Contains(h.cfg.TradingServiceShares, position.ShareName) {
		logrus.WithFields(logrus.Fields{
			"h.cfg.TradingServiceShares": h.cfg.TradingServiceShares,
			"position.ShareName":         position.ShareName,
		}).Errorf("TradingServiceHandler -> OpenPosition -> error: such share doesn't exist")
		return &protocol.OpenPositionResponse{}, fmt.Errorf("TradingServiceHandler -> OpenPosition -> error: such share doesn't exist")
	}

	err = h.s.AddInfoToManager(ctx, position.ProfileID, position.ShareName, position.MoneyAmount, false)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"position.ProfileID":   position.ProfileID,
			"position.ShareName":   position.ShareName,
			"position.MoneyAmount": position.MoneyAmount,
			"backup":               false,
		}).Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}

	pnl, err := h.s.OpenPosition(ctx, &position, false)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"position": position,
			"backup":   false,
		}).Errorf("TradingServiceHandler -> OpenPosition -> %v", err)
		return &protocol.OpenPositionResponse{}, err
	}
	return &protocol.OpenPositionResponse{Pnl: pnl}, nil
}

// ClosePosition closes position by calling method service TradingServiceService ClosePosition
func (h *TradingServiceHandler) ClosePosition(ctx context.Context, req *protocol.ClosePositionRequest) (*protocol.ClosePositionResponse, error) {
	positionID, err := h.validationID(ctx, req.PositionID)
	if err != nil {
		logrus.WithField("req.PositionID", req.PositionID).Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}

	err = h.s.ClosePosition(positionID)
	if err != nil {
		logrus.WithField("positionID", positionID).Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}
	return &protocol.ClosePositionResponse{}, nil
}

// ReadAllOpenedPositionsByProfileID returnes info about all positions that wasn't be closed by exact profile
func (h *TradingServiceHandler) ReadAllOpenedPositionsByProfileID(ctx context.Context, req *protocol.ReadAllOpenedPositionsByProfileIDRequest) (
	*protocol.ReadAllOpenedPositionsByProfileIDResponse, error) {
	profileID, err := h.validationID(ctx, req.ProfileID)
	if err != nil {
		logrus.WithField("req.ProfileID", req.ProfileID).Errorf("TradingServiceHandler -> ReadAllOpenedPositionsByProfileID -> %v", err)
		return &protocol.ReadAllOpenedPositionsByProfileIDResponse{}, err
	}
	openedPositions, err := h.s.ReadAllOpenedPositionsByProfileID(ctx, profileID)
	if err != nil {
		logrus.WithField("profileID", profileID).Errorf("TradingServiceHandler -> ReadAllOpenedPositionsByProfileID -> %v", err)
		return &protocol.ReadAllOpenedPositionsByProfileIDResponse{}, err
	}
	openedPositionsResponse := make([]*protocol.OpenedPosition, 0)
	for _, openedPosition := range openedPositions {
		tempProtoPosition := protocol.OpenedPosition{
			Position: &protocol.Position{
				ProfileID:   openedPosition.ProfileID.String(),
				PositionID:  openedPosition.PositionID.String(),
				TakeProfit:  openedPosition.TakeProfit,
				StopLoss:    openedPosition.StopLoss,
				Vector:      openedPosition.Vector,
				ShareName:   openedPosition.ShareName,
				MoneyAmount: openedPosition.MoneyAmount,
			},
			ShareStartPrice: openedPosition.ShareStartPrice,
			ShareAmount:     openedPosition.ShareAmount,
			OpenedTime:      timestamppb.New(openedPosition.OpenedTime),
		}

		openedPositionsResponse = append(openedPositionsResponse, &tempProtoPosition)
	}
	return &protocol.ReadAllOpenedPositionsByProfileIDResponse{OpenedPositions: openedPositionsResponse}, nil
}
