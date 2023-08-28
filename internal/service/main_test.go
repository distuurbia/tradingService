package service

import (
	"context"
	"os"
	"testing"

	"github.com/distuurbia/tradingService/internal/model"
	"github.com/google/uuid"
)

var (
	testShareApple = model.Share{
		Name:  "Apple",
		Price: 100,
	}
	testShareTesla = model.Share{
		Name:  "Tesla",
		Price: 200,
	}
	testPosition = model.OpenedPosition{
		ShareStartPrice: 100,
		ShareEndPrice:   200,
		ShareAmount:     3,
		Position: model.Position{
			PositionID:  uuid.New(),
			ProfileID:   uuid.New(),
			MoneyAmount: 300,
			StopLoss:    10,
			TakeProfit:  200,
			ShareName:   "Tesla",
			ShortOrLong: "long",
		},
	}
)

const (
	insertActionPayload = `{"action" : "INSERT", "data" : {"positionid":"971fccd8-d211-4ebb-88c7-0049541d9845","profileid":"db85739e-abd7-4a65-9128-7e1528bf0c5a","shortorlong":"short","sharename":"Tesla","shareamount":1.5481699858211602,"sharestartprice":64.59239031620892,"shareendprice":null,"stoploss":1000,"takeprofit":1,"openedtime":"2023-08-28T06:41:55.724626+00:00","closedtime":null}}`
	updateActionPayload = `{"action" : "UPDATE", "data" : {"positionid":"971fccd8-d211-4ebb-88c7-0049541d9845","profileid":"db85739e-abd7-4a65-9128-7e1528bf0c5a","shortorlong":"short","sharename":"Tesla","shareamount":1.5481699858211602,"sharestartprice":64.59239031620892,"shareendprice":60.26545466162589,"stoploss":1000,"takeprofit":1,"openedtime":"2023-08-28T06:41:55.724626+00:00","closedtime":"2023-08-28T06:42:09.811703+00:00"}}`
)

type PriceServiceRepositoryHelper struct {
}

func (h *PriceServiceRepositoryHelper) Subscribe(_ context.Context, _ uuid.UUID, subscribersShares chan model.Share, errSubscribe chan error) {
	subscribersShares <- testShareApple
	subscribersShares <- testShareTesla
	errSubscribe <- nil
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}
