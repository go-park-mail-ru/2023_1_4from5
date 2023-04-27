// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: creator.proto

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

type KeywordMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keyword string `protobuf:"bytes,1,opt,name=Keyword,proto3" json:"Keyword,omitempty"`
}

func (x *KeywordMessage) Reset() {
	*x = KeywordMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeywordMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeywordMessage) ProtoMessage() {}

func (x *KeywordMessage) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeywordMessage.ProtoReflect.Descriptor instead.
func (*KeywordMessage) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{0}
}

func (x *KeywordMessage) GetKeyword() string {
	if x != nil {
		return x.Keyword
	}
	return ""
}

type Creator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	UserID         string `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	CreatorName    string `protobuf:"bytes,3,opt,name=CreatorName,proto3" json:"CreatorName,omitempty"`
	CreatorPhoto   string `protobuf:"bytes,4,opt,name=CreatorPhoto,proto3" json:"CreatorPhoto,omitempty"`
	ProfilePhoto   string `protobuf:"bytes,5,opt,name=ProfilePhoto,proto3" json:"ProfilePhoto,omitempty"`
	FollowersCount int64  `protobuf:"varint,6,opt,name=FollowersCount,proto3" json:"FollowersCount,omitempty"`
	Description    string `protobuf:"bytes,7,opt,name=Description,proto3" json:"Description,omitempty"`
	PostsCount     int64  `protobuf:"varint,8,opt,name=PostsCount,proto3" json:"PostsCount,omitempty"`
}

func (x *Creator) Reset() {
	*x = Creator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Creator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Creator) ProtoMessage() {}

func (x *Creator) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Creator.ProtoReflect.Descriptor instead.
func (*Creator) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{1}
}

func (x *Creator) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Creator) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *Creator) GetCreatorName() string {
	if x != nil {
		return x.CreatorName
	}
	return ""
}

func (x *Creator) GetCreatorPhoto() string {
	if x != nil {
		return x.CreatorPhoto
	}
	return ""
}

func (x *Creator) GetProfilePhoto() string {
	if x != nil {
		return x.ProfilePhoto
	}
	return ""
}

func (x *Creator) GetFollowersCount() int64 {
	if x != nil {
		return x.FollowersCount
	}
	return 0
}

func (x *Creator) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Creator) GetPostsCount() int64 {
	if x != nil {
		return x.PostsCount
	}
	return 0
}

type CreatorsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Creators []*Creator `protobuf:"bytes,9,rep,name=Creators,proto3" json:"Creators,omitempty"`
	Error    string     `protobuf:"bytes,10,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *CreatorsMessage) Reset() {
	*x = CreatorsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatorsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatorsMessage) ProtoMessage() {}

func (x *CreatorsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatorsMessage.ProtoReflect.Descriptor instead.
func (*CreatorsMessage) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{2}
}

func (x *CreatorsMessage) GetCreators() []*Creator {
	if x != nil {
		return x.Creators
	}
	return nil
}

func (x *CreatorsMessage) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UserCreatorMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID    string `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	CreatorID string `protobuf:"bytes,2,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
}

func (x *UserCreatorMessage) Reset() {
	*x = UserCreatorMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserCreatorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserCreatorMessage) ProtoMessage() {}

func (x *UserCreatorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserCreatorMessage.ProtoReflect.Descriptor instead.
func (*UserCreatorMessage) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{3}
}

func (x *UserCreatorMessage) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *UserCreatorMessage) GetCreatorID() string {
	if x != nil {
		return x.CreatorID
	}
	return ""
}

type CreatorPage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatorInfo   *Creator              `protobuf:"bytes,1,opt,name=CreatorInfo,proto3" json:"CreatorInfo,omitempty"`
	AimInfo       *Aim                  `protobuf:"bytes,2,opt,name=AimInfo,proto3" json:"AimInfo,omitempty"`
	IsMyPage      bool                  `protobuf:"varint,3,opt,name=IsMyPage,proto3" json:"IsMyPage,omitempty"`
	Follows       bool                  `protobuf:"varint,4,opt,name=Follows,proto3" json:"Follows,omitempty"`
	Posts         []*Post               `protobuf:"bytes,5,rep,name=Posts,proto3" json:"Posts,omitempty"`
	Subscriptions []*proto.Subscription `protobuf:"bytes,6,rep,name=Subscriptions,proto3" json:"Subscriptions,omitempty"`
	Error         string                `protobuf:"bytes,7,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *CreatorPage) Reset() {
	*x = CreatorPage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatorPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatorPage) ProtoMessage() {}

func (x *CreatorPage) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatorPage.ProtoReflect.Descriptor instead.
func (*CreatorPage) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{4}
}

func (x *CreatorPage) GetCreatorInfo() *Creator {
	if x != nil {
		return x.CreatorInfo
	}
	return nil
}

func (x *CreatorPage) GetAimInfo() *Aim {
	if x != nil {
		return x.AimInfo
	}
	return nil
}

func (x *CreatorPage) GetIsMyPage() bool {
	if x != nil {
		return x.IsMyPage
	}
	return false
}

func (x *CreatorPage) GetFollows() bool {
	if x != nil {
		return x.Follows
	}
	return false
}

func (x *CreatorPage) GetPosts() []*Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

func (x *CreatorPage) GetSubscriptions() []*proto.Subscription {
	if x != nil {
		return x.Subscriptions
	}
	return nil
}

func (x *CreatorPage) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type Aim struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Creator     string `protobuf:"bytes,1,opt,name=Creator,proto3" json:"Creator,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=Description,proto3" json:"Description,omitempty"`
	MoneyNeeded int64  `protobuf:"varint,3,opt,name=MoneyNeeded,proto3" json:"MoneyNeeded,omitempty"`
	MoneyGot    int64  `protobuf:"varint,4,opt,name=MoneyGot,proto3" json:"MoneyGot,omitempty"`
}

func (x *Aim) Reset() {
	*x = Aim{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Aim) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Aim) ProtoMessage() {}

func (x *Aim) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Aim.ProtoReflect.Descriptor instead.
func (*Aim) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{5}
}

func (x *Aim) GetCreator() string {
	if x != nil {
		return x.Creator
	}
	return ""
}

func (x *Aim) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Aim) GetMoneyNeeded() int64 {
	if x != nil {
		return x.MoneyNeeded
	}
	return 0
}

func (x *Aim) GetMoneyGot() int64 {
	if x != nil {
		return x.MoneyGot
	}
	return 0
}

type Post struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string                `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	CreatorID       string                `protobuf:"bytes,2,opt,name=CreatorID,proto3" json:"CreatorID,omitempty"`
	Creation        string                `protobuf:"bytes,3,opt,name=Creation,proto3" json:"Creation,omitempty"`
	LikesCount      int64                 `protobuf:"varint,4,opt,name=LikesCount,proto3" json:"LikesCount,omitempty"`
	ProfilePhoto    string                `protobuf:"bytes,5,opt,name=ProfilePhoto,proto3" json:"ProfilePhoto,omitempty"`
	Title           string                `protobuf:"bytes,6,opt,name=Title,proto3" json:"Title,omitempty"`
	Text            string                `protobuf:"bytes,7,opt,name=Text,proto3" json:"Text,omitempty"`
	IsAvailable     bool                  `protobuf:"varint,8,opt,name=IsAvailable,proto3" json:"IsAvailable,omitempty"`
	IsLiked         bool                  `protobuf:"varint,9,opt,name=IsLiked,proto3" json:"IsLiked,omitempty"`
	PostAttachments []*Attachment         `protobuf:"bytes,10,rep,name=PostAttachments,proto3" json:"PostAttachments,omitempty"`
	Subscriptions   []*proto.Subscription `protobuf:"bytes,11,rep,name=Subscriptions,proto3" json:"Subscriptions,omitempty"`
}

func (x *Post) Reset() {
	*x = Post{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Post) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Post) ProtoMessage() {}

func (x *Post) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Post.ProtoReflect.Descriptor instead.
func (*Post) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{6}
}

func (x *Post) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Post) GetCreatorID() string {
	if x != nil {
		return x.CreatorID
	}
	return ""
}

func (x *Post) GetCreation() string {
	if x != nil {
		return x.Creation
	}
	return ""
}

func (x *Post) GetLikesCount() int64 {
	if x != nil {
		return x.LikesCount
	}
	return 0
}

func (x *Post) GetProfilePhoto() string {
	if x != nil {
		return x.ProfilePhoto
	}
	return ""
}

func (x *Post) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Post) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Post) GetIsAvailable() bool {
	if x != nil {
		return x.IsAvailable
	}
	return false
}

func (x *Post) GetIsLiked() bool {
	if x != nil {
		return x.IsLiked
	}
	return false
}

func (x *Post) GetPostAttachments() []*Attachment {
	if x != nil {
		return x.PostAttachments
	}
	return nil
}

func (x *Post) GetSubscriptions() []*proto.Subscription {
	if x != nil {
		return x.Subscriptions
	}
	return nil
}

type Attachment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID   string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=Type,proto3" json:"Type,omitempty"`
}

func (x *Attachment) Reset() {
	*x = Attachment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_creator_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attachment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attachment) ProtoMessage() {}

func (x *Attachment) ProtoReflect() protoreflect.Message {
	mi := &file_creator_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attachment.ProtoReflect.Descriptor instead.
func (*Attachment) Descriptor() ([]byte, []int) {
	return file_creator_proto_rawDescGZIP(), []int{7}
}

func (x *Attachment) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Attachment) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_creator_proto protoreflect.FileDescriptor

var file_creator_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2a, 0x0a,
	0x0e, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x85, 0x02, 0x0a, 0x07, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x20, 0x0a,
	0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x22, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x50, 0x68,
	0x6f, 0x74, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x68,
	0x6f, 0x74, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x0a, 0x0e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x20, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x22, 0x4d, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x08, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x73,
	0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x52, 0x08, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x22, 0x4a, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1c,
	0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x44, 0x22, 0xfe, 0x01, 0x0a,
	0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x50, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x0b,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x0b, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1e, 0x0a, 0x07, 0x41, 0x69, 0x6d, 0x49,
	0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x41, 0x69, 0x6d, 0x52,
	0x07, 0x41, 0x69, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x73, 0x4d, 0x79,
	0x50, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x49, 0x73, 0x4d, 0x79,
	0x50, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x73, 0x12, 0x1b,
	0x0a, 0x05, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e,
	0x50, 0x6f, 0x73, 0x74, 0x52, 0x05, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x3a, 0x0a, 0x0d, 0x53,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x7f, 0x0a,
	0x03, 0x41, 0x69, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x20,
	0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0b, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x4e, 0x65, 0x65, 0x64, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x4e, 0x65, 0x65, 0x64,
	0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x47, 0x6f, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x47, 0x6f, 0x74, 0x22, 0xed,
	0x02, 0x0a, 0x04, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x6f, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x6f, 0x72, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x68, 0x6f, 0x74,
	0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x50, 0x68, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54,
	0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12,
	0x20, 0x0a, 0x0b, 0x49, 0x73, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x49, 0x73, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x49, 0x73, 0x4c, 0x69, 0x6b, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x49, 0x73, 0x4c, 0x69, 0x6b, 0x65, 0x64, 0x12, 0x35, 0x0a, 0x0f, 0x50,
	0x6f, 0x73, 0x74, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0a,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x12, 0x3a, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0d, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x30,
	0x0a, 0x0a, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x32, 0x75, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x33, 0x0a, 0x0c, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x73, 0x12, 0x0f, 0x2e, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x10, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x50, 0x61,
	0x67, 0x65, 0x12, 0x13, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0c, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x50, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x2f,
	0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_creator_proto_rawDescOnce sync.Once
	file_creator_proto_rawDescData = file_creator_proto_rawDesc
)

func file_creator_proto_rawDescGZIP() []byte {
	file_creator_proto_rawDescOnce.Do(func() {
		file_creator_proto_rawDescData = protoimpl.X.CompressGZIP(file_creator_proto_rawDescData)
	})
	return file_creator_proto_rawDescData
}

var file_creator_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_creator_proto_goTypes = []interface{}{
	(*KeywordMessage)(nil),     // 0: KeywordMessage
	(*Creator)(nil),            // 1: Creator
	(*CreatorsMessage)(nil),    // 2: CreatorsMessage
	(*UserCreatorMessage)(nil), // 3: UserCreatorMessage
	(*CreatorPage)(nil),        // 4: CreatorPage
	(*Aim)(nil),                // 5: Aim
	(*Post)(nil),               // 6: Post
	(*Attachment)(nil),         // 7: Attachment
	(*proto.Subscription)(nil), // 8: common.Subscription
}
var file_creator_proto_depIdxs = []int32{
	1, // 0: CreatorsMessage.Creators:type_name -> Creator
	1, // 1: CreatorPage.CreatorInfo:type_name -> Creator
	5, // 2: CreatorPage.AimInfo:type_name -> Aim
	6, // 3: CreatorPage.Posts:type_name -> Post
	8, // 4: CreatorPage.Subscriptions:type_name -> common.Subscription
	7, // 5: Post.PostAttachments:type_name -> Attachment
	8, // 6: Post.Subscriptions:type_name -> common.Subscription
	0, // 7: CreatorService.FindCreators:input_type -> KeywordMessage
	3, // 8: CreatorService.GetPage:input_type -> UserCreatorMessage
	2, // 9: CreatorService.FindCreators:output_type -> CreatorsMessage
	4, // 10: CreatorService.GetPage:output_type -> CreatorPage
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_creator_proto_init() }
func file_creator_proto_init() {
	if File_creator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_creator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeywordMessage); i {
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
		file_creator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Creator); i {
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
		file_creator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatorsMessage); i {
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
		file_creator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserCreatorMessage); i {
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
		file_creator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatorPage); i {
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
		file_creator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Aim); i {
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
		file_creator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Post); i {
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
		file_creator_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attachment); i {
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
			RawDescriptor: file_creator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_creator_proto_goTypes,
		DependencyIndexes: file_creator_proto_depIdxs,
		MessageInfos:      file_creator_proto_msgTypes,
	}.Build()
	File_creator_proto = out.File
	file_creator_proto_rawDesc = nil
	file_creator_proto_goTypes = nil
	file_creator_proto_depIdxs = nil
}
