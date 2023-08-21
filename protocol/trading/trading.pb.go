// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: trading.proto

package trading

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Position struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MoneyAmount float64 `protobuf:"fixed64,1,opt,name=moneyAmount,proto3" json:"moneyAmount,omitempty"`
	StopLoss    float64 `protobuf:"fixed64,2,opt,name=stopLoss,proto3" json:"stopLoss,omitempty"`
	TakeProfit  float64 `protobuf:"fixed64,3,opt,name=takeProfit,proto3" json:"takeProfit,omitempty"`
	ProfileID   string  `protobuf:"bytes,4,opt,name=profileID,proto3" json:"profileID,omitempty"`
	PositionID  string  `protobuf:"bytes,5,opt,name=positionID,proto3" json:"positionID,omitempty"`
	Vector      string  `protobuf:"bytes,6,opt,name=vector,proto3" json:"vector,omitempty"`
	ShareName   string  `protobuf:"bytes,7,opt,name=shareName,proto3" json:"shareName,omitempty"`
}

func (x *Position) Reset() {
	*x = Position{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Position) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Position) ProtoMessage() {}

func (x *Position) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Position.ProtoReflect.Descriptor instead.
func (*Position) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{0}
}

func (x *Position) GetMoneyAmount() float64 {
	if x != nil {
		return x.MoneyAmount
	}
	return 0
}

func (x *Position) GetStopLoss() float64 {
	if x != nil {
		return x.StopLoss
	}
	return 0
}

func (x *Position) GetTakeProfit() float64 {
	if x != nil {
		return x.TakeProfit
	}
	return 0
}

func (x *Position) GetProfileID() string {
	if x != nil {
		return x.ProfileID
	}
	return ""
}

func (x *Position) GetPositionID() string {
	if x != nil {
		return x.PositionID
	}
	return ""
}

func (x *Position) GetVector() string {
	if x != nil {
		return x.Vector
	}
	return ""
}

func (x *Position) GetShareName() string {
	if x != nil {
		return x.ShareName
	}
	return ""
}

type OpenedPosition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShareStartPrice float64              `protobuf:"fixed64,1,opt,name=shareStartPrice,proto3" json:"shareStartPrice,omitempty"`
	ShareAmount     float64              `protobuf:"fixed64,2,opt,name=shareAmount,proto3" json:"shareAmount,omitempty"`
	OpenedTime      *timestamp.Timestamp `protobuf:"bytes,3,opt,name=openedTime,proto3" json:"openedTime,omitempty"`
	Position        *Position            `protobuf:"bytes,4,opt,name=position,proto3" json:"position,omitempty"`
}

func (x *OpenedPosition) Reset() {
	*x = OpenedPosition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenedPosition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenedPosition) ProtoMessage() {}

func (x *OpenedPosition) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenedPosition.ProtoReflect.Descriptor instead.
func (*OpenedPosition) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{1}
}

func (x *OpenedPosition) GetShareStartPrice() float64 {
	if x != nil {
		return x.ShareStartPrice
	}
	return 0
}

func (x *OpenedPosition) GetShareAmount() float64 {
	if x != nil {
		return x.ShareAmount
	}
	return 0
}

func (x *OpenedPosition) GetOpenedTime() *timestamp.Timestamp {
	if x != nil {
		return x.OpenedTime
	}
	return nil
}

func (x *OpenedPosition) GetPosition() *Position {
	if x != nil {
		return x.Position
	}
	return nil
}

type ReadAllOpenedPositionsByProfileIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProfileID string `protobuf:"bytes,1,opt,name=profileID,proto3" json:"profileID,omitempty"`
}

func (x *ReadAllOpenedPositionsByProfileIDRequest) Reset() {
	*x = ReadAllOpenedPositionsByProfileIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadAllOpenedPositionsByProfileIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadAllOpenedPositionsByProfileIDRequest) ProtoMessage() {}

func (x *ReadAllOpenedPositionsByProfileIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadAllOpenedPositionsByProfileIDRequest.ProtoReflect.Descriptor instead.
func (*ReadAllOpenedPositionsByProfileIDRequest) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{2}
}

func (x *ReadAllOpenedPositionsByProfileIDRequest) GetProfileID() string {
	if x != nil {
		return x.ProfileID
	}
	return ""
}

type ReadAllOpenedPositionsByProfileIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpenedPositions []*OpenedPosition `protobuf:"bytes,1,rep,name=openedPositions,proto3" json:"openedPositions,omitempty"`
}

func (x *ReadAllOpenedPositionsByProfileIDResponse) Reset() {
	*x = ReadAllOpenedPositionsByProfileIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadAllOpenedPositionsByProfileIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadAllOpenedPositionsByProfileIDResponse) ProtoMessage() {}

func (x *ReadAllOpenedPositionsByProfileIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadAllOpenedPositionsByProfileIDResponse.ProtoReflect.Descriptor instead.
func (*ReadAllOpenedPositionsByProfileIDResponse) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{3}
}

func (x *ReadAllOpenedPositionsByProfileIDResponse) GetOpenedPositions() []*OpenedPosition {
	if x != nil {
		return x.OpenedPositions
	}
	return nil
}

type OpenPositionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Position *Position `protobuf:"bytes,1,opt,name=position,proto3" json:"position,omitempty"`
}

func (x *OpenPositionRequest) Reset() {
	*x = OpenPositionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenPositionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenPositionRequest) ProtoMessage() {}

func (x *OpenPositionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenPositionRequest.ProtoReflect.Descriptor instead.
func (*OpenPositionRequest) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{4}
}

func (x *OpenPositionRequest) GetPosition() *Position {
	if x != nil {
		return x.Position
	}
	return nil
}

type OpenPositionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pnl float64 `protobuf:"fixed64,1,opt,name=pnl,proto3" json:"pnl,omitempty"`
}

func (x *OpenPositionResponse) Reset() {
	*x = OpenPositionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpenPositionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpenPositionResponse) ProtoMessage() {}

func (x *OpenPositionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpenPositionResponse.ProtoReflect.Descriptor instead.
func (*OpenPositionResponse) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{5}
}

func (x *OpenPositionResponse) GetPnl() float64 {
	if x != nil {
		return x.Pnl
	}
	return 0
}

type ClosePositionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PositionID string `protobuf:"bytes,1,opt,name=positionID,proto3" json:"positionID,omitempty"`
}

func (x *ClosePositionRequest) Reset() {
	*x = ClosePositionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClosePositionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosePositionRequest) ProtoMessage() {}

func (x *ClosePositionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClosePositionRequest.ProtoReflect.Descriptor instead.
func (*ClosePositionRequest) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{6}
}

func (x *ClosePositionRequest) GetPositionID() string {
	if x != nil {
		return x.PositionID
	}
	return ""
}

type ClosePositionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClosePositionResponse) Reset() {
	*x = ClosePositionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClosePositionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosePositionResponse) ProtoMessage() {}

func (x *ClosePositionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClosePositionResponse.ProtoReflect.Descriptor instead.
func (*ClosePositionResponse) Descriptor() ([]byte, []int) {
	return file_trading_proto_rawDescGZIP(), []int{7}
}

var File_trading_proto protoreflect.FileDescriptor

var file_trading_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xdc, 0x01, 0x0a, 0x08, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a,
	0x0b, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x0b, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x4c, 0x6f, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x4c, 0x6f, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74,
	0x61, 0x6b, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0a, 0x74, 0x61, 0x6b, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x68, 0x61, 0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x61, 0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0xbf, 0x01, 0x0a, 0x0e, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x0f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3a,
	0x0a, 0x0a, 0x6f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a,
	0x6f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x08, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x48, 0x0a, 0x28, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x4f, 0x70, 0x65, 0x6e,
	0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x22, 0x66, 0x0a, 0x29, 0x52,
	0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0f, 0x6f, 0x70, 0x65, 0x6e,
	0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0f, 0x6f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x22, 0x3c, 0x0a, 0x13, 0x4f, 0x70, 0x65, 0x6e, 0x50, 0x6f, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x08, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x28, 0x0a, 0x14, 0x4f, 0x70, 0x65, 0x6e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x6e, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x70, 0x6e, 0x6c, 0x22, 0x36, 0x0a, 0x14, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x22, 0x17, 0x0a, 0x15, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x96, 0x02, 0x0a,
	0x15, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x0c, 0x4f, 0x70, 0x65, 0x6e, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x2e, 0x4f, 0x70, 0x65, 0x6e, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x4f,
	0x70, 0x65, 0x6e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x7c, 0x0a, 0x21, 0x52, 0x65, 0x61, 0x64, 0x41,
	0x6c, 0x6c, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x29, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c,
	0x6c, 0x4f, 0x70, 0x65, 0x6e, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x42, 0x79, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x69, 0x73, 0x74, 0x75, 0x75, 0x72, 0x62, 0x69, 0x61, 0x2f, 0x74,
	0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trading_proto_rawDescOnce sync.Once
	file_trading_proto_rawDescData = file_trading_proto_rawDesc
)

func file_trading_proto_rawDescGZIP() []byte {
	file_trading_proto_rawDescOnce.Do(func() {
		file_trading_proto_rawDescData = protoimpl.X.CompressGZIP(file_trading_proto_rawDescData)
	})
	return file_trading_proto_rawDescData
}

var file_trading_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_trading_proto_goTypes = []interface{}{
	(*Position)(nil),       // 0: Position
	(*OpenedPosition)(nil), // 1: OpenedPosition
	(*ReadAllOpenedPositionsByProfileIDRequest)(nil),  // 2: ReadAllOpenedPositionsByProfileIDRequest
	(*ReadAllOpenedPositionsByProfileIDResponse)(nil), // 3: ReadAllOpenedPositionsByProfileIDResponse
	(*OpenPositionRequest)(nil),                       // 4: OpenPositionRequest
	(*OpenPositionResponse)(nil),                      // 5: OpenPositionResponse
	(*ClosePositionRequest)(nil),                      // 6: ClosePositionRequest
	(*ClosePositionResponse)(nil),                     // 7: ClosePositionResponse
	(*timestamp.Timestamp)(nil),                       // 8: google.protobuf.Timestamp
}
var file_trading_proto_depIdxs = []int32{
	8, // 0: OpenedPosition.openedTime:type_name -> google.protobuf.Timestamp
	0, // 1: OpenedPosition.position:type_name -> Position
	1, // 2: ReadAllOpenedPositionsByProfileIDResponse.openedPositions:type_name -> OpenedPosition
	0, // 3: OpenPositionRequest.position:type_name -> Position
	4, // 4: TradingServiceService.OpenPosition:input_type -> OpenPositionRequest
	6, // 5: TradingServiceService.ClosePosition:input_type -> ClosePositionRequest
	2, // 6: TradingServiceService.ReadAllOpenedPositionsByProfileID:input_type -> ReadAllOpenedPositionsByProfileIDRequest
	5, // 7: TradingServiceService.OpenPosition:output_type -> OpenPositionResponse
	7, // 8: TradingServiceService.ClosePosition:output_type -> ClosePositionResponse
	3, // 9: TradingServiceService.ReadAllOpenedPositionsByProfileID:output_type -> ReadAllOpenedPositionsByProfileIDResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_trading_proto_init() }
func file_trading_proto_init() {
	if File_trading_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trading_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Position); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenedPosition); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadAllOpenedPositionsByProfileIDRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadAllOpenedPositionsByProfileIDResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenPositionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpenPositionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClosePositionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trading_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClosePositionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_trading_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trading_proto_goTypes,
		DependencyIndexes: file_trading_proto_depIdxs,
		MessageInfos:      file_trading_proto_msgTypes,
	}.Build()
	File_trading_proto = out.File
	file_trading_proto_rawDesc = nil
	file_trading_proto_goTypes = nil
	file_trading_proto_depIdxs = nil
}
