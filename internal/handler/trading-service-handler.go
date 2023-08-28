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

// TradingServiceService is an interface of the service TradingServiceService
type TradingServiceService interface {
	AddPosition(ctx context.Context, position *model.Position) error
	ClosePosition(ctx context.Context, profileID, positionID uuid.UUID) (float64, error)
	ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error)
}

// TradingServiceHandler implements the  TradingServiceService interface also contains an object of type *validator.Validate and *config.Config
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

// validationID validates given id and parses it to uuid.UUID type
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

// AddPosition opens position by calling method AddPosition of service TradingServiceService
func (h *TradingServiceHandler) AddPosition(ctx context.Context, req *protocol.AddPositionRequest) (*protocol.AddPositionResponse, error) {
	profileID, err := h.validationID(ctx, req.Position.ProfileID)
	if err != nil {
		logrus.WithField("req.Position.ProfileID", req.Position.ProfileID).Errorf("TradingServiceHandler -> AddPosition -> %v", err)
		return &protocol.AddPositionResponse{}, err
	}
	positionID, err := h.validationID(ctx, req.Position.PositionID)
	if err != nil {
		logrus.WithField("req.Position.PositionID", req.Position.PositionID).Errorf("TradingServiceHandler -> AddPosition -> %v", err)
		return &protocol.AddPositionResponse{}, err
	}
	position := model.Position{
		ProfileID:   profileID,
		PositionID:  positionID,
		StopLoss:    req.Position.StopLoss,
		TakeProfit:  req.Position.TakeProfit,
		MoneyAmount: req.Position.MoneyAmount,
		ShareName:   req.Position.ShareName,
		ShortOrLong: req.Position.ShortOrLong,
	}
	err = h.validate.StructCtx(ctx, position)
	if err != nil {
		logrus.WithField("position", position).Errorf("TradingServiceHandler -> AddPosition -> %v", err)
		return &protocol.AddPositionResponse{}, err
	}

	if !strings.Contains(h.cfg.TradingServiceShares, position.ShareName) {
		logrus.WithFields(logrus.Fields{
			"h.cfg.TradingServiceShares": h.cfg.TradingServiceShares,
			"position.ShareName":         position.ShareName,
		}).Errorf("TradingServiceHandler -> AddPosition -> error: such share doesn't exist")
		return &protocol.AddPositionResponse{}, fmt.Errorf("TradingServiceHandler -> AddPosition -> error: such share doesn't exist")
	}

	err = h.s.AddPosition(ctx, &position)
	if err != nil {
		logrus.WithField("position", position).Errorf("TradingServiceHandler -> AddPosition -> %v", err)
		return &protocol.AddPositionResponse{}, err
	}

	return &protocol.AddPositionResponse{}, nil
}

// ClosePosition closes position by calling method service TradingServiceService ClosePosition
func (h *TradingServiceHandler) ClosePosition(ctx context.Context, req *protocol.ClosePositionRequest) (*protocol.ClosePositionResponse, error) {
	profileID, err := h.validationID(ctx, req.ProfileID)
	if err != nil {
		logrus.WithField("req.PositionID", req.PositionID).Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}

	positionID, err := h.validationID(ctx, req.PositionID)
	if err != nil {
		logrus.WithField("req.PositionID", req.PositionID).Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}

	pnl, err := h.s.ClosePosition(ctx, profileID, positionID)
	if err != nil {
		logrus.WithField("positionID", positionID).Errorf("TradingServiceHandler -> ClosePosition -> %v", err)
		return &protocol.ClosePositionResponse{}, err
	}
	return &protocol.ClosePositionResponse{Pnl: pnl}, nil
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
				ShortOrLong: openedPosition.ShortOrLong,
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
