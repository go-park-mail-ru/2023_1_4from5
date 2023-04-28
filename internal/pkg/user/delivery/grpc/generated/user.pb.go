// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: user.proto

package generated

import (
	proto "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
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

type FollowMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID    string `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	CreatorID string `protobuf:"bytes,2,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
}

func (x *FollowMessage) Reset() {
	*x = FollowMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowMessage) ProtoMessage() {}

func (x *FollowMessage) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowMessage.ProtoReflect.Descriptor instead.
func (*FollowMessage) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *FollowMessage) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *FollowMessage) GetCreatorID() string {
	if x != nil {
		return x.CreatorID
	}
	return ""
}

type SubscriptionDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	UserID     string `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	MonthCount int64  `protobuf:"varint,3,opt,name=MonthCount,proto3" json:"MonthCount,omitempty"`
	Money      int64  `protobuf:"varint,4,opt,name=Money,proto3" json:"Money,omitempty"`
	CreatorID  string `protobuf:"bytes,5,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
}

func (x *SubscriptionDetails) Reset() {
	*x = SubscriptionDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscriptionDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscriptionDetails) ProtoMessage() {}

func (x *SubscriptionDetails) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscriptionDetails.ProtoReflect.Descriptor instead.
func (*SubscriptionDetails) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *SubscriptionDetails) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SubscriptionDetails) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *SubscriptionDetails) GetMonthCount() int64 {
	if x != nil {
		return x.MonthCount
	}
	return 0
}

func (x *SubscriptionDetails) GetMoney() int64 {
	if x != nil {
		return x.Money
	}
	return 0
}

func (x *SubscriptionDetails) GetCreatorID() string {
	if x != nil {
		return x.CreatorID
	}
	return ""
}

type ImageID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *ImageID) Reset() {
	*x = ImageID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageID) ProtoMessage() {}

func (x *ImageID) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageID.ProtoReflect.Descriptor instead.
func (*ImageID) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *ImageID) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *ImageID) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UserProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login        string `protobuf:"bytes,1,opt,name=Login,proto3" json:"Login,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	ProfilePhoto string `protobuf:"bytes,3,opt,name=ProfilePhoto,proto3" json:"ProfilePhoto,omitempty"`
	Registration string `protobuf:"bytes,4,opt,name=Registration,proto3" json:"Registration,omitempty"`
	Error        string `protobuf:"bytes,5,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *UserProfile) Reset() {
	*x = UserProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserProfile) ProtoMessage() {}

func (x *UserProfile) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserProfile.ProtoReflect.Descriptor instead.
func (*UserProfile) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *UserProfile) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *UserProfile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserProfile) GetProfilePhoto() string {
	if x != nil {
		return x.ProfilePhoto
	}
	return ""
}

func (x *UserProfile) GetRegistration() string {
	if x != nil {
		return x.Registration
	}
	return ""
}

func (x *UserProfile) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdatePasswordMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID   string `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
}

func (x *UpdatePasswordMessage) Reset() {
	*x = UpdatePasswordMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePasswordMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePasswordMessage) ProtoMessage() {}

func (x *UpdatePasswordMessage) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePasswordMessage.ProtoReflect.Descriptor instead.
func (*UpdatePasswordMessage) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{4}
}

func (x *UpdatePasswordMessage) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *UpdatePasswordMessage) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type UpdateProfileInfoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login  string `protobuf:"bytes,1,opt,name=Login,proto3" json:"Login,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	UserID string `protobuf:"bytes,3,opt,name=UserID,proto3" json:"UserID,omitempty"`
}

func (x *UpdateProfileInfoMessage) Reset() {
	*x = UpdateProfileInfoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProfileInfoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProfileInfoMessage) ProtoMessage() {}

func (x *UpdateProfileInfoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProfileInfoMessage.ProtoReflect.Descriptor instead.
func (*UpdateProfileInfoMessage) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateProfileInfoMessage) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *UpdateProfileInfoMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateProfileInfoMessage) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

type DonateMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatorID  string `protobuf:"bytes,1,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
	MoneyCount int64  `protobuf:"varint,2,opt,name=MoneyCount,proto3" json:"MoneyCount,omitempty"`
	UserID     string `protobuf:"bytes,3,opt,name=UserID,proto3" json:"UserID,omitempty"`
}

func (x *DonateMessage) Reset() {
	*x = DonateMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DonateMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DonateMessage) ProtoMessage() {}

func (x *DonateMessage) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DonateMessage.ProtoReflect.Descriptor instead.
func (*DonateMessage) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{6}
}

func (x *DonateMessage) GetCreatorID() string {
	if x != nil {
		return x.CreatorID
	}
	return ""
}

func (x *DonateMessage) GetMoneyCount() int64 {
	if x != nil {
		return x.MoneyCount
	}
	return 0
}

func (x *DonateMessage) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

type DonateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MoneyCount int64  `protobuf:"varint,1,opt,name=MoneyCount,proto3" json:"MoneyCount,omitempty"`
	Error      string `protobuf:"bytes,2,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *DonateResponse) Reset() {
	*x = DonateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DonateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DonateResponse) ProtoMessage() {}

func (x *DonateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DonateResponse.ProtoReflect.Descriptor instead.
func (*DonateResponse) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{7}
}

func (x *DonateResponse) GetMoneyCount() int64 {
	if x != nil {
		return x.MoneyCount
	}
	return 0
}

func (x *DonateResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type BecameCreatorInfoMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=Description,proto3" json:"Description,omitempty"`
	UserID      string `protobuf:"bytes,3,opt,name=UserID,proto3" json:"UserID,omitempty"`
}

func (x *BecameCreatorInfoMessage) Reset() {
	*x = BecameCreatorInfoMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BecameCreatorInfoMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BecameCreatorInfoMessage) ProtoMessage() {}

func (x *BecameCreatorInfoMessage) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BecameCreatorInfoMessage.ProtoReflect.Descriptor instead.
func (*BecameCreatorInfoMessage) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{8}
}

func (x *BecameCreatorInfoMessage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BecameCreatorInfoMessage) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *BecameCreatorInfoMessage) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

type SubscriptionsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subscriptions []*proto.Subscription `protobuf:"bytes,8,rep,name=Subscriptions,proto3" json:"Subscriptions,omitempty"`
	Error         string                `protobuf:"bytes,9,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *SubscriptionsMessage) Reset() {
	*x = SubscriptionsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscriptionsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscriptionsMessage) ProtoMessage() {}

func (x *SubscriptionsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_user_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscriptionsMessage.ProtoReflect.Descriptor instead.
func (*SubscriptionsMessage) Descriptor() ([]byte, []int) {
	return file_user_proto_rawDescGZIP(), []int{9}
}

func (x *SubscriptionsMessage) GetSubscriptions() []*proto.Subscription {
	if x != nil {
		return x.Subscriptions
	}
	return nil
}

func (x *SubscriptionsMessage) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x0d, 0x46, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x44,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49,
	0x44, 0x22, 0x91, 0x01, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x6f, 0x72, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x6f, 0x72, 0x49, 0x44, 0x22, 0x35, 0x0a, 0x07, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x44,
	0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x95, 0x01, 0x0a,
	0x0b, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14,
	0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x4b, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x5c, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22,
	0x65, 0x0a, 0x0d, 0x44, 0x6f, 0x6e, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x12, 0x1e,
	0x0a, 0x0a, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x46, 0x0a, 0x0e, 0x44, 0x6f, 0x6e, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x6f, 0x6e, 0x65,
	0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x4d, 0x6f,
	0x6e, 0x65, 0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x68,
	0x0a, 0x18, 0x42, 0x65, 0x63, 0x61, 0x6d, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x68, 0x0a, 0x14, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x3a, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x32, 0xac, 0x04, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x0e, 0x2e, 0x46,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2b, 0x0a,
	0x08, 0x55, 0x6e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x0e, 0x2e, 0x46, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x09, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x14, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x1a, 0x0d, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x31,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x13, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x1a, 0x0c, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x22,
	0x00, 0x12, 0x2e, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f,
	0x12, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x08, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x49, 0x44, 0x22,
	0x00, 0x12, 0x39, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x16, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x11,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x19, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2b, 0x0a,
	0x06, 0x44, 0x6f, 0x6e, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x2e, 0x44, 0x6f, 0x6e, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0f, 0x2e, 0x44, 0x6f, 0x6e, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x0d, 0x42, 0x65,
	0x63, 0x6f, 0x6d, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x19, 0x2e, 0x42, 0x65,
	0x63, 0x61, 0x6d, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x55, 0x55, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x41,
	0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x55, 0x49,
	0x44, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x15, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x42, 0x2b, 0x5a, 0x29, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_proto_rawDescOnce sync.Once
	file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
	file_user_proto_rawDescOnce.Do(func() {
		file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
	})
	return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_user_proto_goTypes = []interface{}{
	(*FollowMessage)(nil),            // 0: FollowMessage
	(*SubscriptionDetails)(nil),      // 1: SubscriptionDetails
	(*ImageID)(nil),                  // 2: ImageID
	(*UserProfile)(nil),              // 3: UserProfile
	(*UpdatePasswordMessage)(nil),    // 4: UpdatePasswordMessage
	(*UpdateProfileInfoMessage)(nil), // 5: UpdateProfileInfoMessage
	(*DonateMessage)(nil),            // 6: DonateMessage
	(*DonateResponse)(nil),           // 7: DonateResponse
	(*BecameCreatorInfoMessage)(nil), // 8: BecameCreatorInfoMessage
	(*SubscriptionsMessage)(nil),     // 9: SubscriptionsMessage
	(*proto.Subscription)(nil),       // 10: common.Subscription
	(*proto.UUIDMessage)(nil),        // 11: common.UUIDMessage
	(*proto.Empty)(nil),              // 12: common.Empty
	(*proto.UUIDResponse)(nil),       // 13: common.UUIDResponse
}
var file_user_proto_depIdxs = []int32{
	10, // 0: SubscriptionsMessage.Subscriptions:type_name -> common.Subscription
	0,  // 1: UserService.Follow:input_type -> FollowMessage
	0,  // 2: UserService.Unfollow:input_type -> FollowMessage
	1,  // 3: UserService.Subscribe:input_type -> SubscriptionDetails
	11, // 4: UserService.GetProfile:input_type -> common.UUIDMessage
	11, // 5: UserService.UpdatePhoto:input_type -> common.UUIDMessage
	4,  // 6: UserService.UpdatePassword:input_type -> UpdatePasswordMessage
	5,  // 7: UserService.UpdateProfileInfo:input_type -> UpdateProfileInfoMessage
	6,  // 8: UserService.Donate:input_type -> DonateMessage
	8,  // 9: UserService.BecomeCreator:input_type -> BecameCreatorInfoMessage
	11, // 10: UserService.UserSubscriptions:input_type -> common.UUIDMessage
	12, // 11: UserService.Follow:output_type -> common.Empty
	12, // 12: UserService.Unfollow:output_type -> common.Empty
	12, // 13: UserService.Subscribe:output_type -> common.Empty
	3,  // 14: UserService.GetProfile:output_type -> UserProfile
	2,  // 15: UserService.UpdatePhoto:output_type -> ImageID
	12, // 16: UserService.UpdatePassword:output_type -> common.Empty
	12, // 17: UserService.UpdateProfileInfo:output_type -> common.Empty
	7,  // 18: UserService.Donate:output_type -> DonateResponse
	13, // 19: UserService.BecomeCreator:output_type -> common.UUIDResponse
	9,  // 20: UserService.UserSubscriptions:output_type -> SubscriptionsMessage
	11, // [11:21] is the sub-list for method output_type
	1,  // [1:11] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
	if File_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowMessage); i {
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
		file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscriptionDetails); i {
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
		file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageID); i {
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
		file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserProfile); i {
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
		file_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePasswordMessage); i {
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
		file_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProfileInfoMessage); i {
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
		file_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DonateMessage); i {
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
		file_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DonateResponse); i {
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
		file_user_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BecameCreatorInfoMessage); i {
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
		file_user_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscriptionsMessage); i {
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
			RawDescriptor: file_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_proto_goTypes,
		DependencyIndexes: file_user_proto_depIdxs,
		MessageInfos:      file_user_proto_msgTypes,
	}.Build()
	File_user_proto = out.File
	file_user_proto_rawDesc = nil
	file_user_proto_goTypes = nil
	file_user_proto_depIdxs = nil
}
