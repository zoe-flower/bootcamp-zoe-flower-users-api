// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: justeat/menucard-v1.proto

package justeat

import (
	_ "github.com/flypay/events/pkg/flyt"
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

type MenuCardCreatedV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Define your proto fields here
	MenuCardId   int32  `protobuf:"varint,1,opt,name=menu_card_id,json=MenuCardId,proto3" json:"menu_card_id,omitempty"`
	MenuCardType string `protobuf:"bytes,2,opt,name=menu_card_type,json=MenuCardType,proto3" json:"menu_card_type,omitempty"`
	RestaurantId int32  `protobuf:"varint,4,opt,name=restaurant_id,json=RestaurantId,proto3" json:"restaurant_id,omitempty"`
	Tenant       string `protobuf:"bytes,5,opt,name=tenant,json=Tenant,proto3" json:"tenant,omitempty"`
}

func (x *MenuCardCreatedV1) Reset() {
	*x = MenuCardCreatedV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_justeat_menucard_v1_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MenuCardCreatedV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MenuCardCreatedV1) ProtoMessage() {}

func (x *MenuCardCreatedV1) ProtoReflect() protoreflect.Message {
	mi := &file_justeat_menucard_v1_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MenuCardCreatedV1.ProtoReflect.Descriptor instead.
func (*MenuCardCreatedV1) Descriptor() ([]byte, []int) {
	return file_justeat_menucard_v1_proto_rawDescGZIP(), []int{0}
}

func (x *MenuCardCreatedV1) GetMenuCardId() int32 {
	if x != nil {
		return x.MenuCardId
	}
	return 0
}

func (x *MenuCardCreatedV1) GetMenuCardType() string {
	if x != nil {
		return x.MenuCardType
	}
	return ""
}

func (x *MenuCardCreatedV1) GetRestaurantId() int32 {
	if x != nil {
		return x.RestaurantId
	}
	return 0
}

func (x *MenuCardCreatedV1) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

type MenuCardUpdatedV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Define your proto fields here
	MenuCardId   int32  `protobuf:"varint,1,opt,name=menu_card_id,json=MenuCardId,proto3" json:"menu_card_id,omitempty"`
	MenuCardType string `protobuf:"bytes,2,opt,name=menu_card_type,json=MenuCardType,proto3" json:"menu_card_type,omitempty"`
	RestaurantId int32  `protobuf:"varint,4,opt,name=restaurant_id,json=RestaurantId,proto3" json:"restaurant_id,omitempty"`
	Tenant       string `protobuf:"bytes,5,opt,name=tenant,json=Tenant,proto3" json:"tenant,omitempty"`
}

func (x *MenuCardUpdatedV1) Reset() {
	*x = MenuCardUpdatedV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_justeat_menucard_v1_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MenuCardUpdatedV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MenuCardUpdatedV1) ProtoMessage() {}

func (x *MenuCardUpdatedV1) ProtoReflect() protoreflect.Message {
	mi := &file_justeat_menucard_v1_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MenuCardUpdatedV1.ProtoReflect.Descriptor instead.
func (*MenuCardUpdatedV1) Descriptor() ([]byte, []int) {
	return file_justeat_menucard_v1_proto_rawDescGZIP(), []int{1}
}

func (x *MenuCardUpdatedV1) GetMenuCardId() int32 {
	if x != nil {
		return x.MenuCardId
	}
	return 0
}

func (x *MenuCardUpdatedV1) GetMenuCardType() string {
	if x != nil {
		return x.MenuCardType
	}
	return ""
}

func (x *MenuCardUpdatedV1) GetRestaurantId() int32 {
	if x != nil {
		return x.RestaurantId
	}
	return 0
}

func (x *MenuCardUpdatedV1) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

type MenuCardDeletedV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Define your proto fields here
	MenuCardId   int32  `protobuf:"varint,1,opt,name=menu_card_id,json=MenuCardId,proto3" json:"menu_card_id,omitempty"`
	MenuCardType string `protobuf:"bytes,2,opt,name=menu_card_type,json=MenuCardType,proto3" json:"menu_card_type,omitempty"`
	RestaurantId int32  `protobuf:"varint,4,opt,name=restaurant_id,json=RestaurantId,proto3" json:"restaurant_id,omitempty"`
	Tenant       string `protobuf:"bytes,5,opt,name=tenant,json=Tenant,proto3" json:"tenant,omitempty"`
}

func (x *MenuCardDeletedV1) Reset() {
	*x = MenuCardDeletedV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_justeat_menucard_v1_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MenuCardDeletedV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MenuCardDeletedV1) ProtoMessage() {}

func (x *MenuCardDeletedV1) ProtoReflect() protoreflect.Message {
	mi := &file_justeat_menucard_v1_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MenuCardDeletedV1.ProtoReflect.Descriptor instead.
func (*MenuCardDeletedV1) Descriptor() ([]byte, []int) {
	return file_justeat_menucard_v1_proto_rawDescGZIP(), []int{2}
}

func (x *MenuCardDeletedV1) GetMenuCardId() int32 {
	if x != nil {
		return x.MenuCardId
	}
	return 0
}

func (x *MenuCardDeletedV1) GetMenuCardType() string {
	if x != nil {
		return x.MenuCardType
	}
	return ""
}

func (x *MenuCardDeletedV1) GetRestaurantId() int32 {
	if x != nil {
		return x.RestaurantId
	}
	return 0
}

func (x *MenuCardDeletedV1) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

var File_justeat_menucard_v1_proto protoreflect.FileDescriptor

var file_justeat_menucard_v1_proto_rawDesc = []byte{
	0x0a, 0x19, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2f, 0x6d, 0x65, 0x6e, 0x75, 0x63, 0x61,
	0x72, 0x64, 0x2d, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6a, 0x75, 0x73,
	0x74, 0x65, 0x61, 0x74, 0x1a, 0x15, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6a, 0x75, 0x73,
	0x74, 0x65, 0x61, 0x74, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xfc, 0x01, 0x0a, 0x11, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x56, 0x31, 0x12, 0x20, 0x0a, 0x0c, 0x6d, 0x65, 0x6e, 0x75,
	0x5f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x65,
	0x6e, 0x75, 0x5f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72,
	0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x3a, 0x62, 0x82,
	0xb5, 0x18, 0x1b, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x63,
	0x61, 0x72, 0x64, 0x2e, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x76, 0x31, 0xa2, 0xbb,
	0x18, 0x0f, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0xaa, 0xbb, 0x18, 0x02, 0x75, 0x6b, 0xaa, 0xbb, 0x18, 0x02, 0x65, 0x73, 0xaa, 0xbb, 0x18,
	0x02, 0x69, 0x65, 0xaa, 0xbb, 0x18, 0x02, 0x69, 0x74, 0xaa, 0xbb, 0x18, 0x02, 0x64, 0x6b, 0xaa,
	0xbb, 0x18, 0x02, 0x6e, 0x6f, 0xaa, 0xbb, 0x18, 0x02, 0x61, 0x75, 0xaa, 0xbb, 0x18, 0x02, 0x6e,
	0x7a, 0x22, 0xfc, 0x01, 0x0a, 0x11, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x56, 0x31, 0x12, 0x20, 0x0a, 0x0c, 0x6d, 0x65, 0x6e, 0x75, 0x5f,
	0x63, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x4d,
	0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x65, 0x6e,
	0x75, 0x5f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x3a, 0x62, 0x82, 0xb5,
	0x18, 0x1b, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x63, 0x61,
	0x72, 0x64, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x76, 0x31, 0xa2, 0xbb, 0x18,
	0x0f, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0xaa, 0xbb, 0x18, 0x02, 0x75, 0x6b, 0xaa, 0xbb, 0x18, 0x02, 0x65, 0x73, 0xaa, 0xbb, 0x18, 0x02,
	0x69, 0x65, 0xaa, 0xbb, 0x18, 0x02, 0x69, 0x74, 0xaa, 0xbb, 0x18, 0x02, 0x64, 0x6b, 0xaa, 0xbb,
	0x18, 0x02, 0x6e, 0x6f, 0xaa, 0xbb, 0x18, 0x02, 0x61, 0x75, 0xaa, 0xbb, 0x18, 0x02, 0x6e, 0x7a,
	0x22, 0xfc, 0x01, 0x0a, 0x11, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x56, 0x31, 0x12, 0x20, 0x0a, 0x0c, 0x6d, 0x65, 0x6e, 0x75, 0x5f, 0x63,
	0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x4d, 0x65,
	0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x6d, 0x65, 0x6e, 0x75,
	0x5f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x3a, 0x62, 0x82, 0xb5, 0x18,
	0x1b, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x63, 0x61, 0x72,
	0x64, 0x2e, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x2e, 0x76, 0x31, 0xa2, 0xbb, 0x18, 0x0f,
	0x4d, 0x65, 0x6e, 0x75, 0x43, 0x61, 0x72, 0x64, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0xaa,
	0xbb, 0x18, 0x02, 0x75, 0x6b, 0xaa, 0xbb, 0x18, 0x02, 0x65, 0x73, 0xaa, 0xbb, 0x18, 0x02, 0x69,
	0x65, 0xaa, 0xbb, 0x18, 0x02, 0x69, 0x74, 0xaa, 0xbb, 0x18, 0x02, 0x64, 0x6b, 0xaa, 0xbb, 0x18,
	0x02, 0x6e, 0x6f, 0xaa, 0xbb, 0x18, 0x02, 0x61, 0x75, 0xaa, 0xbb, 0x18, 0x02, 0x6e, 0x7a, 0x42,
	0x80, 0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x42,
	0x0f, 0x4d, 0x65, 0x6e, 0x75, 0x63, 0x61, 0x72, 0x64, 0x56, 0x31, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66,
	0x6c, 0x79, 0x70, 0x61, 0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0xa2, 0x02, 0x03, 0x4a, 0x58, 0x58, 0xaa, 0x02,
	0x07, 0x4a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0xca, 0x02, 0x07, 0x4a, 0x75, 0x73, 0x74, 0x65,
	0x61, 0x74, 0xe2, 0x02, 0x13, 0x4a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x07, 0x4a, 0x75, 0x73, 0x74, 0x65,
	0x61, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_justeat_menucard_v1_proto_rawDescOnce sync.Once
	file_justeat_menucard_v1_proto_rawDescData = file_justeat_menucard_v1_proto_rawDesc
)

func file_justeat_menucard_v1_proto_rawDescGZIP() []byte {
	file_justeat_menucard_v1_proto_rawDescOnce.Do(func() {
		file_justeat_menucard_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_justeat_menucard_v1_proto_rawDescData)
	})
	return file_justeat_menucard_v1_proto_rawDescData
}

var file_justeat_menucard_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_justeat_menucard_v1_proto_goTypes = []interface{}{
	(*MenuCardCreatedV1)(nil), // 0: justeat.MenuCardCreatedV1
	(*MenuCardUpdatedV1)(nil), // 1: justeat.MenuCardUpdatedV1
	(*MenuCardDeletedV1)(nil), // 2: justeat.MenuCardDeletedV1
}
var file_justeat_menucard_v1_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_justeat_menucard_v1_proto_init() }
func file_justeat_menucard_v1_proto_init() {
	if File_justeat_menucard_v1_proto != nil {
		return
	}
	file_justeat_options_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_justeat_menucard_v1_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MenuCardCreatedV1); i {
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
		file_justeat_menucard_v1_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MenuCardUpdatedV1); i {
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
		file_justeat_menucard_v1_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MenuCardDeletedV1); i {
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
			RawDescriptor: file_justeat_menucard_v1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_justeat_menucard_v1_proto_goTypes,
		DependencyIndexes: file_justeat_menucard_v1_proto_depIdxs,
		MessageInfos:      file_justeat_menucard_v1_proto_msgTypes,
	}.Build()
	File_justeat_menucard_v1_proto = out.File
	file_justeat_menucard_v1_proto_rawDesc = nil
	file_justeat_menucard_v1_proto_goTypes = nil
	file_justeat_menucard_v1_proto_depIdxs = nil
}
