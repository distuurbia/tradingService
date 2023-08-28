// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	trading "github.com/distuurbia/tradingService/protocol/trading"
)

// TradingServiceServiceClient is an autogenerated mock type for the TradingServiceServiceClient type
type TradingServiceServiceClient struct {
	mock.Mock
}

// AddPosition provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceServiceClient) AddPosition(ctx context.Context, in *trading.AddPositionRequest, opts ...grpc.CallOption) (*trading.AddPositionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *trading.AddPositionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *trading.AddPositionRequest, ...grpc.CallOption) (*trading.AddPositionResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *trading.AddPositionRequest, ...grpc.CallOption) *trading.AddPositionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*trading.AddPositionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *trading.AddPositionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClosePosition provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceServiceClient) ClosePosition(ctx context.Context, in *trading.ClosePositionRequest, opts ...grpc.CallOption) (*trading.ClosePositionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *trading.ClosePositionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *trading.ClosePositionRequest, ...grpc.CallOption) (*trading.ClosePositionResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *trading.ClosePositionRequest, ...grpc.CallOption) *trading.ClosePositionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*trading.ClosePositionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *trading.ClosePositionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadAllOpenedPositionsByProfileID provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceServiceClient) ReadAllOpenedPositionsByProfileID(ctx context.Context, in *trading.ReadAllOpenedPositionsByProfileIDRequest, opts ...grpc.CallOption) (*trading.ReadAllOpenedPositionsByProfileIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *trading.ReadAllOpenedPositionsByProfileIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *trading.ReadAllOpenedPositionsByProfileIDRequest, ...grpc.CallOption) (*trading.ReadAllOpenedPositionsByProfileIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *trading.ReadAllOpenedPositionsByProfileIDRequest, ...grpc.CallOption) *trading.ReadAllOpenedPositionsByProfileIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*trading.ReadAllOpenedPositionsByProfileIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *trading.ReadAllOpenedPositionsByProfileIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTradingServiceServiceClient creates a new instance of TradingServiceServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTradingServiceServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *TradingServiceServiceClient {
	mock := &TradingServiceServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
