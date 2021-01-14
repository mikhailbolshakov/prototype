// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: proto/mm/mm.proto

package mm

import (
	proto "github.com/golang/protobuf/proto"
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

type CreateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Email    string `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{0}
}

func (x *CreateUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *CreateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type CreateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *CreateUserResponse) Reset() {
	*x = CreateUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserResponse) ProtoMessage() {}

func (x *CreateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserResponse.ProtoReflect.Descriptor instead.
func (*CreateUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{1}
}

func (x *CreateUserResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UserStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status   string `protobuf:"bytes,1,opt,name=Status,proto3" json:"Status,omitempty"`
	MMUserId string `protobuf:"bytes,2,opt,name=MMUserId,proto3" json:"MMUserId,omitempty"`
}

func (x *UserStatus) Reset() {
	*x = UserStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserStatus) ProtoMessage() {}

func (x *UserStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserStatus.ProtoReflect.Descriptor instead.
func (*UserStatus) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{2}
}

func (x *UserStatus) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *UserStatus) GetMMUserId() string {
	if x != nil {
		return x.MMUserId
	}
	return ""
}

type GetUsersStatusesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MMUserIds []string `protobuf:"bytes,1,rep,name=MMUserIds,proto3" json:"MMUserIds,omitempty"`
}

func (x *GetUsersStatusesRequest) Reset() {
	*x = GetUsersStatusesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersStatusesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersStatusesRequest) ProtoMessage() {}

func (x *GetUsersStatusesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersStatusesRequest.ProtoReflect.Descriptor instead.
func (*GetUsersStatusesRequest) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{3}
}

func (x *GetUsersStatusesRequest) GetMMUserIds() []string {
	if x != nil {
		return x.MMUserIds
	}
	return nil
}

type GetUserStatusesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Statuses []*UserStatus `protobuf:"bytes,1,rep,name=Statuses,proto3" json:"Statuses,omitempty"`
}

func (x *GetUserStatusesResponse) Reset() {
	*x = GetUserStatusesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserStatusesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserStatusesResponse) ProtoMessage() {}

func (x *GetUserStatusesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserStatusesResponse.ProtoReflect.Descriptor instead.
func (*GetUserStatusesResponse) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserStatusesResponse) GetStatuses() []*UserStatus {
	if x != nil {
		return x.Statuses
	}
	return nil
}

type CreateClientChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientUserId string   `protobuf:"bytes,1,opt,name=ClientUserId,proto3" json:"ClientUserId,omitempty"`
	Name         string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	DisplayName  string   `protobuf:"bytes,3,opt,name=DisplayName,proto3" json:"DisplayName,omitempty"`
	Subscribers  []string `protobuf:"bytes,4,rep,name=Subscribers,proto3" json:"Subscribers,omitempty"`
}

func (x *CreateClientChannelRequest) Reset() {
	*x = CreateClientChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateClientChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateClientChannelRequest) ProtoMessage() {}

func (x *CreateClientChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateClientChannelRequest.ProtoReflect.Descriptor instead.
func (*CreateClientChannelRequest) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{5}
}

func (x *CreateClientChannelRequest) GetClientUserId() string {
	if x != nil {
		return x.ClientUserId
	}
	return ""
}

func (x *CreateClientChannelRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateClientChannelRequest) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *CreateClientChannelRequest) GetSubscribers() []string {
	if x != nil {
		return x.Subscribers
	}
	return nil
}

type CreateClientChannelResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,2,opt,name=ChannelId,proto3" json:"ChannelId,omitempty"`
}

func (x *CreateClientChannelResponse) Reset() {
	*x = CreateClientChannelResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mm_mm_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateClientChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateClientChannelResponse) ProtoMessage() {}

func (x *CreateClientChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mm_mm_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateClientChannelResponse.ProtoReflect.Descriptor instead.
func (*CreateClientChannelResponse) Descriptor() ([]byte, []int) {
	return file_proto_mm_mm_proto_rawDescGZIP(), []int{6}
}

func (x *CreateClientChannelResponse) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

var File_proto_mm_mm_proto protoreflect.FileDescriptor

var file_proto_mm_mm_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6d, 0x2f, 0x6d, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x02, 0x6d, 0x6d, 0x22, 0x45, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x24,
	0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x64, 0x22, 0x40, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x4d,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4d, 0x4d,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x37, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x4d, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x4d, 0x4d, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73, 0x22,
	0x45, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x08, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6d,
	0x6d, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x08, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x22, 0x98, 0x01, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72,
	0x73, 0x22, 0x3b, 0x0a, 0x1b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x32, 0x96,
	0x01, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x3d, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x15, 0x2e, 0x6d, 0x6d, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x6d, 0x6d, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0x1b, 0x2e, 0x6d, 0x6d,
	0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x6d, 0x6d, 0x2e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x64, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x73, 0x12, 0x58, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1e, 0x2e, 0x6d, 0x6d, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6d, 0x6d, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x07, 0x0a,
	0x05, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x6d, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_mm_mm_proto_rawDescOnce sync.Once
	file_proto_mm_mm_proto_rawDescData = file_proto_mm_mm_proto_rawDesc
)

func file_proto_mm_mm_proto_rawDescGZIP() []byte {
	file_proto_mm_mm_proto_rawDescOnce.Do(func() {
		file_proto_mm_mm_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_mm_mm_proto_rawDescData)
	})
	return file_proto_mm_mm_proto_rawDescData
}

var file_proto_mm_mm_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_mm_mm_proto_goTypes = []interface{}{
	(*CreateUserRequest)(nil),           // 0: mm.CreateUserRequest
	(*CreateUserResponse)(nil),          // 1: mm.CreateUserResponse
	(*UserStatus)(nil),                  // 2: mm.UserStatus
	(*GetUsersStatusesRequest)(nil),     // 3: mm.GetUsersStatusesRequest
	(*GetUserStatusesResponse)(nil),     // 4: mm.GetUserStatusesResponse
	(*CreateClientChannelRequest)(nil),  // 5: mm.CreateClientChannelRequest
	(*CreateClientChannelResponse)(nil), // 6: mm.CreateClientChannelResponse
}
var file_proto_mm_mm_proto_depIdxs = []int32{
	2, // 0: mm.GetUserStatusesResponse.Statuses:type_name -> mm.UserStatus
	0, // 1: mm.Users.CreateUser:input_type -> mm.CreateUserRequest
	3, // 2: mm.Users.GetUsersStatuses:input_type -> mm.GetUsersStatusesRequest
	5, // 3: mm.Channels.CreateClientChannel:input_type -> mm.CreateClientChannelRequest
	1, // 4: mm.Users.CreateUser:output_type -> mm.CreateUserResponse
	4, // 5: mm.Users.GetUsersStatuses:output_type -> mm.GetUserStatusesResponse
	6, // 6: mm.Channels.CreateClientChannel:output_type -> mm.CreateClientChannelResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_mm_mm_proto_init() }
func file_proto_mm_mm_proto_init() {
	if File_proto_mm_mm_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_mm_mm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserRequest); i {
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
		file_proto_mm_mm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserResponse); i {
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
		file_proto_mm_mm_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserStatus); i {
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
		file_proto_mm_mm_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUsersStatusesRequest); i {
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
		file_proto_mm_mm_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserStatusesResponse); i {
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
		file_proto_mm_mm_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateClientChannelRequest); i {
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
		file_proto_mm_mm_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateClientChannelResponse); i {
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
			RawDescriptor: file_proto_mm_mm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_proto_mm_mm_proto_goTypes,
		DependencyIndexes: file_proto_mm_mm_proto_depIdxs,
		MessageInfos:      file_proto_mm_mm_proto_msgTypes,
	}.Build()
	File_proto_mm_mm_proto = out.File
	file_proto_mm_mm_proto_rawDesc = nil
	file_proto_mm_mm_proto_goTypes = nil
	file_proto_mm_mm_proto_depIdxs = nil
}
