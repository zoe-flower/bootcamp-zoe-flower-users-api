// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: flyt/legacy-menus-published.proto

package flyt

import (
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

type MenuType int32

const (
	MenuType_ORDER_AHEAD       MenuType = 0
	MenuType_EAT_IN            MenuType = 1
	MenuType_ORDER_AND_DELIVER MenuType = 2
)

// Enum value maps for MenuType.
var (
	MenuType_name = map[int32]string{
		0: "ORDER_AHEAD",
		1: "EAT_IN",
		2: "ORDER_AND_DELIVER",
	}
	MenuType_value = map[string]int32{
		"ORDER_AHEAD":       0,
		"EAT_IN":            1,
		"ORDER_AND_DELIVER": 2,
	}
)

func (x MenuType) Enum() *MenuType {
	p := new(MenuType)
	*p = x
	return p
}

func (x MenuType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MenuType) Descriptor() protoreflect.EnumDescriptor {
	return file_flyt_legacy_menus_published_proto_enumTypes[0].Descriptor()
}

func (MenuType) Type() protoreflect.EnumType {
	return &file_flyt_legacy_menus_published_proto_enumTypes[0]
}

func (x MenuType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MenuType.Descriptor instead.
func (MenuType) EnumDescriptor() ([]byte, []int) {
	return file_flyt_legacy_menus_published_proto_rawDescGZIP(), []int{0}
}

// Skip message - candidate for being removed when we no longer need to
// publish menus via core platform
type LegacyMenusPublished struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Restaurant *Ident `protobuf:"bytes,1,opt,name=restaurant,proto3" json:"restaurant,omitempty"`
	// Deprecated: Do not use.
	Types       []MenuType                         `protobuf:"varint,2,rep,packed,name=types,proto3,enum=flyt.MenuType" json:"types,omitempty"`
	LegacyMenus []*LegacyMenusPublished_LegacyMenu `protobuf:"bytes,3,rep,name=legacy_menus,json=legacyMenus,proto3" json:"legacy_menus,omitempty"`
}

func (x *LegacyMenusPublished) Reset() {
	*x = LegacyMenusPublished{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_legacy_menus_published_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LegacyMenusPublished) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LegacyMenusPublished) ProtoMessage() {}

func (x *LegacyMenusPublished) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_legacy_menus_published_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LegacyMenusPublished.ProtoReflect.Descriptor instead.
func (*LegacyMenusPublished) Descriptor() ([]byte, []int) {
	return file_flyt_legacy_menus_published_proto_rawDescGZIP(), []int{0}
}

func (x *LegacyMenusPublished) GetRestaurant() *Ident {
	if x != nil {
		return x.Restaurant
	}
	return nil
}

// Deprecated: Do not use.
func (x *LegacyMenusPublished) GetTypes() []MenuType {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *LegacyMenusPublished) GetLegacyMenus() []*LegacyMenusPublished_LegacyMenu {
	if x != nil {
		return x.LegacyMenus
	}
	return nil
}

// Contains the menu type and the url to the payload
type LegacyMenusPublished_LegacyMenu struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type       MenuType `protobuf:"varint,1,opt,name=type,proto3,enum=flyt.MenuType" json:"type,omitempty"`
	PayloadUrl string   `protobuf:"bytes,2,opt,name=payload_url,json=payloadUrl,proto3" json:"payload_url,omitempty"`
}

func (x *LegacyMenusPublished_LegacyMenu) Reset() {
	*x = LegacyMenusPublished_LegacyMenu{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_legacy_menus_published_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LegacyMenusPublished_LegacyMenu) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LegacyMenusPublished_LegacyMenu) ProtoMessage() {}

func (x *LegacyMenusPublished_LegacyMenu) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_legacy_menus_published_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LegacyMenusPublished_LegacyMenu.ProtoReflect.Descriptor instead.
func (*LegacyMenusPublished_LegacyMenu) Descriptor() ([]byte, []int) {
	return file_flyt_legacy_menus_published_proto_rawDescGZIP(), []int{0, 0}
}

func (x *LegacyMenusPublished_LegacyMenu) GetType() MenuType {
	if x != nil {
		return x.Type
	}
	return MenuType_ORDER_AHEAD
}

func (x *LegacyMenusPublished_LegacyMenu) GetPayloadUrl() string {
	if x != nil {
		return x.PayloadUrl
	}
	return ""
}

var File_flyt_legacy_menus_published_proto protoreflect.FileDescriptor

var file_flyt_legacy_menus_published_proto_rawDesc = []byte{
	0x0a, 0x21, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x2d, 0x6d, 0x65,
	0x6e, 0x75, 0x73, 0x2d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66, 0x6c, 0x79, 0x74, 0x1a, 0x15, 0x66, 0x6c, 0x79, 0x74, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x10, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xa6, 0x02, 0x0a, 0x14, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d, 0x65, 0x6e,
	0x75, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x12, 0x2b, 0x0a, 0x0a, 0x72,
	0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x72, 0x65,
	0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x4d,
	0x65, 0x6e, 0x75, 0x54, 0x79, 0x70, 0x65, 0x42, 0x02, 0x18, 0x01, 0x52, 0x05, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x12, 0x48, 0x0a, 0x0c, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x5f, 0x6d, 0x65, 0x6e,
	0x75, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e,
	0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d, 0x65, 0x6e, 0x75, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x65, 0x64, 0x2e, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d, 0x65, 0x6e, 0x75, 0x52,
	0x0b, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d, 0x65, 0x6e, 0x75, 0x73, 0x1a, 0x51, 0x0a, 0x0a,
	0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d, 0x65, 0x6e, 0x75, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e,
	0x4d, 0x65, 0x6e, 0x75, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x55, 0x72, 0x6c, 0x3a,
	0x1a, 0x82, 0xb5, 0x18, 0x16, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x2e, 0x6d, 0x65, 0x6e, 0x75,
	0x73, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x2a, 0x3e, 0x0a, 0x08, 0x4d,
	0x65, 0x6e, 0x75, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x4f, 0x52, 0x44, 0x45, 0x52,
	0x5f, 0x41, 0x48, 0x45, 0x41, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x41, 0x54, 0x5f,
	0x49, 0x4e, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x41, 0x4e,
	0x44, 0x5f, 0x44, 0x45, 0x4c, 0x49, 0x56, 0x45, 0x52, 0x10, 0x02, 0x42, 0x78, 0x0a, 0x08, 0x63,
	0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x42, 0x19, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d,
	0x65, 0x6e, 0x75, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x6c, 0x79, 0x70, 0x61, 0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x66, 0x6c, 0x79, 0x74, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x04,
	0x46, 0x6c, 0x79, 0x74, 0xca, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0xe2, 0x02, 0x10, 0x46, 0x6c,
	0x79, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x04, 0x46, 0x6c, 0x79, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flyt_legacy_menus_published_proto_rawDescOnce sync.Once
	file_flyt_legacy_menus_published_proto_rawDescData = file_flyt_legacy_menus_published_proto_rawDesc
)

func file_flyt_legacy_menus_published_proto_rawDescGZIP() []byte {
	file_flyt_legacy_menus_published_proto_rawDescOnce.Do(func() {
		file_flyt_legacy_menus_published_proto_rawDescData = protoimpl.X.CompressGZIP(file_flyt_legacy_menus_published_proto_rawDescData)
	})
	return file_flyt_legacy_menus_published_proto_rawDescData
}

var file_flyt_legacy_menus_published_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_flyt_legacy_menus_published_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_flyt_legacy_menus_published_proto_goTypes = []interface{}{
	(MenuType)(0),                           // 0: flyt.MenuType
	(*LegacyMenusPublished)(nil),            // 1: flyt.LegacyMenusPublished
	(*LegacyMenusPublished_LegacyMenu)(nil), // 2: flyt.LegacyMenusPublished.LegacyMenu
	(*Ident)(nil),                           // 3: flyt.Ident
}
var file_flyt_legacy_menus_published_proto_depIdxs = []int32{
	3, // 0: flyt.LegacyMenusPublished.restaurant:type_name -> flyt.Ident
	0, // 1: flyt.LegacyMenusPublished.types:type_name -> flyt.MenuType
	2, // 2: flyt.LegacyMenusPublished.legacy_menus:type_name -> flyt.LegacyMenusPublished.LegacyMenu
	0, // 3: flyt.LegacyMenusPublished.LegacyMenu.type:type_name -> flyt.MenuType
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_flyt_legacy_menus_published_proto_init() }
func file_flyt_legacy_menus_published_proto_init() {
	if File_flyt_legacy_menus_published_proto != nil {
		return
	}
	file_flyt_descriptor_proto_init()
	file_flyt_ident_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_flyt_legacy_menus_published_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LegacyMenusPublished); i {
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
		file_flyt_legacy_menus_published_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LegacyMenusPublished_LegacyMenu); i {
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
			RawDescriptor: file_flyt_legacy_menus_published_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flyt_legacy_menus_published_proto_goTypes,
		DependencyIndexes: file_flyt_legacy_menus_published_proto_depIdxs,
		EnumInfos:         file_flyt_legacy_menus_published_proto_enumTypes,
		MessageInfos:      file_flyt_legacy_menus_published_proto_msgTypes,
	}.Build()
	File_flyt_legacy_menus_published_proto = out.File
	file_flyt_legacy_menus_published_proto_rawDesc = nil
	file_flyt_legacy_menus_published_proto_goTypes = nil
	file_flyt_legacy_menus_published_proto_depIdxs = nil
}
