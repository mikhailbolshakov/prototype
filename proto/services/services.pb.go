// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: proto/services/services.proto

package services

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ChangeServicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	ServiceTypeId string `protobuf:"bytes,2,opt,name=ServiceTypeId,proto3" json:"ServiceTypeId,omitempty"`
	Quantity      int32  `protobuf:"varint,3,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
}

func (x *ChangeServicesRequest) Reset() {
	*x = ChangeServicesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_services_services_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeServicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeServicesRequest) ProtoMessage() {}

func (x *ChangeServicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_services_services_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeServicesRequest.ProtoReflect.Descriptor instead.
func (*ChangeServicesRequest) Descriptor() ([]byte, []int) {
	return file_proto_services_services_proto_rawDescGZIP(), []int{0}
}

func (x *ChangeServicesRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ChangeServicesRequest) GetServiceTypeId() string {
	if x != nil {
		return x.ServiceTypeId
	}
	return ""
}

func (x *ChangeServicesRequest) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type Balance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Available int32 `protobuf:"varint,1,opt,name=Available,proto3" json:"Available,omitempty"`
	Delivered int32 `protobuf:"varint,2,opt,name=Delivered,proto3" json:"Delivered,omitempty"`
	Locked    int32 `protobuf:"varint,3,opt,name=Locked,proto3" json:"Locked,omitempty"`
	Total     int32 `protobuf:"varint,4,opt,name=Total,proto3" json:"Total,omitempty"`
}

func (x *Balance) Reset() {
	*x = Balance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_services_services_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Balance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Balance) ProtoMessage() {}

func (x *Balance) ProtoReflect() protoreflect.Message {
	mi := &file_proto_services_services_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Balance.ProtoReflect.Descriptor instead.
func (*Balance) Descriptor() ([]byte, []int) {
	return file_proto_services_services_proto_rawDescGZIP(), []int{1}
}

func (x *Balance) GetAvailable() int32 {
	if x != nil {
		return x.Available
	}
	return 0
}

func (x *Balance) GetDelivered() int32 {
	if x != nil {
		return x.Delivered
	}
	return 0
}

func (x *Balance) GetLocked() int32 {
	if x != nil {
		return x.Locked
	}
	return 0
}

func (x *Balance) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type UserBalance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  string              `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Balance map[string]*Balance `protobuf:"bytes,2,rep,name=Balance,proto3" json:"Balance,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UserBalance) Reset() {
	*x = UserBalance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_services_services_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserBalance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserBalance) ProtoMessage() {}

func (x *UserBalance) ProtoReflect() protoreflect.Message {
	mi := &file_proto_services_services_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserBalance.ProtoReflect.Descriptor instead.
func (*UserBalance) Descriptor() ([]byte, []int) {
	return file_proto_services_services_proto_rawDescGZIP(), []int{2}
}

func (x *UserBalance) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserBalance) GetBalance() map[string]*Balance {
	if x != nil {
		return x.Balance
	}
	return nil
}

type GetBalanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string               `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	At     *timestamp.Timestamp `protobuf:"bytes,2,opt,name=At,proto3" json:"At,omitempty"`
}

func (x *GetBalanceRequest) Reset() {
	*x = GetBalanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_services_services_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBalanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBalanceRequest) ProtoMessage() {}

func (x *GetBalanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_services_services_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBalanceRequest.ProtoReflect.Descriptor instead.
func (*GetBalanceRequest) Descriptor() ([]byte, []int) {
	return file_proto_services_services_proto_rawDescGZIP(), []int{3}
}

func (x *GetBalanceRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetBalanceRequest) GetAt() *timestamp.Timestamp {
	if x != nil {
		return x.At
	}
	return nil
}

type DeliveryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        string `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	ServiceTypeId string `protobuf:"bytes,2,opt,name=ServiceTypeId,proto3" json:"ServiceTypeId,omitempty"`
	Details       []byte `protobuf:"bytes,3,opt,name=Details,proto3" json:"Details,omitempty"`
}

func (x *DeliveryRequest) Reset() {
	*x = DeliveryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_services_services_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliveryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliveryRequest) ProtoMessage() {}

func (x *DeliveryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_services_services_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliveryRequest.ProtoReflect.Descriptor instead.
func (*DeliveryRequest) Descriptor() ([]byte, []int) {
	return file_proto_services_services_proto_rawDescGZIP(), []int{4}
}

func (x *DeliveryRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DeliveryRequest) GetServiceTypeId() string {
	if x != nil {
		return x.ServiceTypeId
	}
	return ""
}

func (x *DeliveryRequest) GetDetails() []byte {
	if x != nil {
		return x.Details
	}
	return nil
}

type Delivery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string               `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	UserId        string               `protobuf:"bytes,2,opt,name=UserId,proto3" json:"UserId,omitempty"`
	ServiceTypeId string               `protobuf:"bytes,3,opt,name=ServiceTypeId,proto3" json:"ServiceTypeId,omitempty"`
	Status        string               `protobuf:"bytes,4,opt,name=Status,proto3" json:"Status,omitempty"`
	StartTime     *timestamp.Timestamp `protobuf:"bytes,5,opt,name=StartTime,proto3" json:"StartTime,omitempty"`
	FinishTime    *timestamp.Timestamp `protobuf:"bytes,6,opt,name=FinishTime,proto3" json:"FinishTime,omitempty"`
	Details       []byte               `protobuf:"bytes,7,opt,name=Details,proto3" json:"Details,omitempty"`
}

func (x *Delivery) Reset() {
	*x = Delivery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_services_services_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Delivery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Delivery) ProtoMessage() {}

func (x *Delivery) ProtoReflect() protoreflect.Message {
	mi := &file_proto_services_services_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Delivery.ProtoReflect.Descriptor instead.
func (*Delivery) Descriptor() ([]byte, []int) {
	return file_proto_services_services_proto_rawDescGZIP(), []int{5}
}

func (x *Delivery) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Delivery) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Delivery) GetServiceTypeId() string {
	if x != nil {
		return x.ServiceTypeId
	}
	return ""
}

func (x *Delivery) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Delivery) GetStartTime() *timestamp.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *Delivery) GetFinishTime() *timestamp.Timestamp {
	if x != nil {
		return x.FinishTime
	}
	return nil
}

func (x *Delivery) GetDetails() []byte {
	if x != nil {
		return x.Details
	}
	return nil
}

var File_proto_services_services_proto protoreflect.FileDescriptor

var file_proto_services_services_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x71, 0x0a, 0x15, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x73, 0x0a,
	0x07, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x41, 0x76, 0x61,
	0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65,
	0x72, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x44, 0x65, 0x6c, 0x69, 0x76,
	0x65, 0x72, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x54, 0x6f, 0x74,
	0x61, 0x6c, 0x22, 0xb2, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x42, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3c, 0x0a, 0x07, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x07, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x1a, 0x4d, 0x0a, 0x0c, 0x42, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x57, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x02, 0x41, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x41, 0x74,
	0x22, 0x69, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x80, 0x02, 0x0a, 0x08,
	0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x24, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x38,
	0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x0a, 0x46, 0x69, 0x6e, 0x69,
	0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x32, 0xdf,
	0x02, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12,
	0x3f, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x00,
	0x12, 0x42, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1b,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x6c,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x08, 0x57, 0x72, 0x69, 0x74, 0x65, 0x4f, 0x66, 0x66,
	0x12, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x04, 0x4c, 0x6f,
	0x63, 0x6b, 0x12, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x0f,
	0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x19, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76,
	0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x22, 0x00,
	0x42, 0x10, 0x5a, 0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_services_services_proto_rawDescOnce sync.Once
	file_proto_services_services_proto_rawDescData = file_proto_services_services_proto_rawDesc
)

func file_proto_services_services_proto_rawDescGZIP() []byte {
	file_proto_services_services_proto_rawDescOnce.Do(func() {
		file_proto_services_services_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_services_services_proto_rawDescData)
	})
	return file_proto_services_services_proto_rawDescData
}

var file_proto_services_services_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_services_services_proto_goTypes = []interface{}{
	(*ChangeServicesRequest)(nil), // 0: services.ChangeServicesRequest
	(*Balance)(nil),               // 1: services.Balance
	(*UserBalance)(nil),           // 2: services.UserBalance
	(*GetBalanceRequest)(nil),     // 3: services.GetBalanceRequest
	(*DeliveryRequest)(nil),       // 4: services.DeliveryRequest
	(*Delivery)(nil),              // 5: services.Delivery
	nil,                           // 6: services.UserBalance.BalanceEntry
	(*timestamp.Timestamp)(nil),   // 7: google.protobuf.Timestamp
}
var file_proto_services_services_proto_depIdxs = []int32{
	6,  // 0: services.UserBalance.Balance:type_name -> services.UserBalance.BalanceEntry
	7,  // 1: services.GetBalanceRequest.At:type_name -> google.protobuf.Timestamp
	7,  // 2: services.Delivery.StartTime:type_name -> google.protobuf.Timestamp
	7,  // 3: services.Delivery.FinishTime:type_name -> google.protobuf.Timestamp
	1,  // 4: services.UserBalance.BalanceEntry.value:type_name -> services.Balance
	0,  // 5: services.UserServices.Add:input_type -> services.ChangeServicesRequest
	3,  // 6: services.UserServices.GetBalance:input_type -> services.GetBalanceRequest
	0,  // 7: services.UserServices.WriteOff:input_type -> services.ChangeServicesRequest
	0,  // 8: services.UserServices.Lock:input_type -> services.ChangeServicesRequest
	4,  // 9: services.UserServices.DeliveryService:input_type -> services.DeliveryRequest
	2,  // 10: services.UserServices.Add:output_type -> services.UserBalance
	2,  // 11: services.UserServices.GetBalance:output_type -> services.UserBalance
	2,  // 12: services.UserServices.WriteOff:output_type -> services.UserBalance
	2,  // 13: services.UserServices.Lock:output_type -> services.UserBalance
	5,  // 14: services.UserServices.DeliveryService:output_type -> services.Delivery
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_proto_services_services_proto_init() }
func file_proto_services_services_proto_init() {
	if File_proto_services_services_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_services_services_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeServicesRequest); i {
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
		file_proto_services_services_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Balance); i {
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
		file_proto_services_services_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserBalance); i {
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
		file_proto_services_services_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBalanceRequest); i {
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
		file_proto_services_services_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliveryRequest); i {
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
		file_proto_services_services_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Delivery); i {
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
			RawDescriptor: file_proto_services_services_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_services_services_proto_goTypes,
		DependencyIndexes: file_proto_services_services_proto_depIdxs,
		MessageInfos:      file_proto_services_services_proto_msgTypes,
	}.Build()
	File_proto_services_services_proto = out.File
	file_proto_services_services_proto_rawDesc = nil
	file_proto_services_services_proto_goTypes = nil
	file_proto_services_services_proto_depIdxs = nil
}
