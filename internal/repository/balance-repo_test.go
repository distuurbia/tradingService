package repository

import (
	"context"
	"testing"

	balanceProtocol "github.com/distuurbia/balance/protocol/balance"
	"github.com/distuurbia/balance/protocol/balance/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWithdrawOnPosition(t *testing.T) {
	client := new(mocks.BalanceServiceClient)
	client.On("GetBalance", mock.Anything, mock.Anything).Return(&balanceProtocol.GetBalanceResponse{TotalBalance: 500}, nil)
	client.On("AddBalanceChange", mock.Anything, mock.Anything).Return(&balanceProtocol.AddBalanceChangeResponse{}, nil)
	r := NewBalanceRepository(client)

	err := r.WithdrawOnPosition(context.Background(), uuid.New().String(), 300)
	require.NoError(t, err)
}

func TestMoneyBack(t *testing.T) {
	client := new(mocks.BalanceServiceClient)
	client.On("AddBalanceChange", mock.Anything, mock.Anything).Return(&balanceProtocol.AddBalanceChangeResponse{}, nil)
	r := NewBalanceRepository(client)

	err := r.MoneyBackWithPnl(context.Background(), uuid.New().String(), 300)
	require.NoError(t, err)
}
