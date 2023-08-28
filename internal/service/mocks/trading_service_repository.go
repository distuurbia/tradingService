// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/distuurbia/tradingService/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// TradingServiceRepository is an autogenerated mock type for the TradingServiceRepository type
type TradingServiceRepository struct {
	mock.Mock
}

// AddPosition provides a mock function with given fields: ctx, position, shareAmount, shareStartPrice
func (_m *TradingServiceRepository) AddPosition(ctx context.Context, position *model.Position, shareAmount float64, shareStartPrice float64) error {
	ret := _m.Called(ctx, position, shareAmount, shareStartPrice)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Position, float64, float64) error); ok {
		r0 = rf(ctx, position, shareAmount, shareStartPrice)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BackupAllOpenedPositions provides a mock function with given fields: ctx
func (_m *TradingServiceRepository) BackupAllOpenedPositions(ctx context.Context) ([]*model.OpenedPosition, error) {
	ret := _m.Called(ctx)

	var r0 []*model.OpenedPosition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.OpenedPosition, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.OpenedPosition); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OpenedPosition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClosePosition provides a mock function with given fields: ctx, positionID, shareEndPrice
func (_m *TradingServiceRepository) ClosePosition(ctx context.Context, positionID uuid.UUID, shareEndPrice float64) error {
	ret := _m.Called(ctx, positionID, shareEndPrice)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, float64) error); ok {
		r0 = rf(ctx, positionID, shareEndPrice)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadAllOpenedPositionsByProfileID provides a mock function with given fields: ctx, profileID
func (_m *TradingServiceRepository) ReadAllOpenedPositionsByProfileID(ctx context.Context, profileID uuid.UUID) ([]*model.OpenedPosition, error) {
	ret := _m.Called(ctx, profileID)

	var r0 []*model.OpenedPosition
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]*model.OpenedPosition, error)); ok {
		return rf(ctx, profileID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*model.OpenedPosition); ok {
		r0 = rf(ctx, profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OpenedPosition)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, profileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadPositionRow provides a mock function with given fields: ctx, positionID
func (_m *TradingServiceRepository) ReadPositionRow(ctx context.Context, positionID uuid.UUID) (float64, float64, float64, string, error) {
	ret := _m.Called(ctx, positionID)

	var r0 float64
	var r1 float64
	var r2 float64
	var r3 string
	var r4 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (float64, float64, float64, string, error)); ok {
		return rf(ctx, positionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) float64); ok {
		r0 = rf(ctx, positionID)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) float64); ok {
		r1 = rf(ctx, positionID)
	} else {
		r1 = ret.Get(1).(float64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uuid.UUID) float64); ok {
		r2 = rf(ctx, positionID)
	} else {
		r2 = ret.Get(2).(float64)
	}

	if rf, ok := ret.Get(3).(func(context.Context, uuid.UUID) string); ok {
		r3 = rf(ctx, positionID)
	} else {
		r3 = ret.Get(3).(string)
	}

	if rf, ok := ret.Get(4).(func(context.Context, uuid.UUID) error); ok {
		r4 = rf(ctx, positionID)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}

// ReadShareNameByPositionID provides a mock function with given fields: ctx, positionID
func (_m *TradingServiceRepository) ReadShareNameByPositionID(ctx context.Context, positionID uuid.UUID) (string, error) {
	ret := _m.Called(ctx, positionID)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (string, error)); ok {
		return rf(ctx, positionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) string); ok {
		r0 = rf(ctx, positionID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, positionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTradingServiceRepository creates a new instance of TradingServiceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTradingServiceRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TradingServiceRepository {
	mock := &TradingServiceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
