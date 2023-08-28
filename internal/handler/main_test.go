package handler

import (
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/distuurbia/tradingService/internal/config"
	"github.com/distuurbia/tradingService/internal/model"
	"github.com/distuurbia/tradingService/protocol/trading"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	validate          *validator.Validate
	cfg               config.Config
	testProtoPosition = &trading.Position{
		ProfileID:   uuid.New().String(),
		PositionID:  uuid.New().String(),
		MoneyAmount: 300,
		StopLoss:    10,
		TakeProfit:  200,
		ShareName:   "Apple",
		ShortOrLong: "long",
	}
	testOpenedPosition = model.OpenedPosition{
		Position: model.Position{
			ProfileID:   uuid.New(),
			PositionID:  uuid.New(),
			MoneyAmount: 300,
			StopLoss:    10,
			TakeProfit:  200,
			ShareName:   "Apple",
			ShortOrLong: "long",
		},
	}
)

func TestMain(m *testing.M) {
	validate = validator.New()
	if errCfg := env.Parse(&cfg); errCfg != nil {
		logrus.Fatalf("main -> %v", errCfg)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}
