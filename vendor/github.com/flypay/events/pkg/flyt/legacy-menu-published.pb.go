// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: flyt/legacy-menu-published.proto

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

// Skip message - candidate for being removed when we no longer need to
// publish menus via core platform
type LegacyMenuPublished struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Restaurant *Ident   `protobuf:"bytes,1,opt,name=restaurant,proto3" json:"restaurant,omitempty"`
	Type       MenuType `protobuf:"varint,2,opt,name=type,proto3,enum=flyt.MenuType" json:"type,omitempty"`
	// MenuID is deprecated because we do not need to send this to skip any longer
	//
	// Deprecated: Do not use.
	MenuId     int64  `protobuf:"varint,3,opt,name=menu_id,json=menuId,proto3" json:"menu_id,omitempty"`
	PayloadUrl string `protobuf:"bytes,4,opt,name=payload_url,json=payloadUrl,proto3" json:"payload_url,omitempty"`
}

func (x *LegacyMenuPublished) Reset() {
	*x = LegacyMenuPublished{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_legacy_menu_published_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LegacyMenuPublished) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LegacyMenuPublished) ProtoMessage() {}

func (x *LegacyMenuPublished) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_legacy_menu_published_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LegacyMenuPublished.ProtoReflect.Descriptor instead.
func (*LegacyMenuPublished) Descriptor() ([]byte, []int) {
	return file_flyt_legacy_menu_published_proto_rawDescGZIP(), []int{0}
}

func (x *LegacyMenuPublished) GetRestaurant() *Ident {
	if x != nil {
		return x.Restaurant
	}
	return nil
}

func (x *LegacyMenuPublished) GetType() MenuType {
	if x != nil {
		return x.Type
	}
	return MenuType_ORDER_AHEAD
}

// Deprecated: Do not use.
func (x *LegacyMenuPublished) GetMenuId() int64 {
	if x != nil {
		return x.MenuId
	}
	return 0
}

func (x *LegacyMenuPublished) GetPayloadUrl() string {
	if x != nil {
		return x.PayloadUrl
	}
	return ""
}

var File_flyt_legacy_menu_published_proto protoreflect.FileDescriptor

var file_flyt_legacy_menu_published_proto_rawDesc = []byte{
	0x0a, 0x20, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x2d, 0x6d, 0x65,
	0x6e, 0x75, 0x2d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x66, 0x6c, 0x79, 0x74, 0x1a, 0x15, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x10, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x21, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x2d, 0x6d,
	0x65, 0x6e, 0x75, 0x73, 0x2d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbf, 0x01, 0x0a, 0x13, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d,
	0x65, 0x6e, 0x75, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x12, 0x2b, 0x0a, 0x0a,
	0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x72,
	0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x4d,
	0x65, 0x6e, 0x75, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a,
	0x07, 0x6d, 0x65, 0x6e, 0x75, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02,
	0x18, 0x01, 0x52, 0x06, 0x6d, 0x65, 0x6e, 0x75, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x55, 0x72, 0x6c, 0x3a, 0x19, 0x82, 0xb5, 0x18,
	0x15, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x2e, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x42, 0x77, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c,
	0x79, 0x74, 0x42, 0x18, 0x4c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x4d, 0x65, 0x6e, 0x75, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x21,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x79, 0x70, 0x61,
	0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x66, 0x6c, 0x79,
	0x74, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0xca, 0x02,
	0x04, 0x46, 0x6c, 0x79, 0x74, 0xe2, 0x02, 0x10, 0x46, 0x6c, 0x79, 0x74, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flyt_legacy_menu_published_proto_rawDescOnce sync.Once
	file_flyt_legacy_menu_published_proto_rawDescData = file_flyt_legacy_menu_published_proto_rawDesc
)

func file_flyt_legacy_menu_published_proto_rawDescGZIP() []byte {
	file_flyt_legacy_menu_published_proto_rawDescOnce.Do(func() {
		file_flyt_legacy_menu_published_proto_rawDescData = protoimpl.X.CompressGZIP(file_flyt_legacy_menu_published_proto_rawDescData)
	})
	return file_flyt_legacy_menu_published_proto_rawDescData
}

var file_flyt_legacy_menu_published_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_flyt_legacy_menu_published_proto_goTypes = []interface{}{
	(*LegacyMenuPublished)(nil), // 0: flyt.LegacyMenuPublished
	(*Ident)(nil),               // 1: flyt.Ident
	(MenuType)(0),               // 2: flyt.MenuType
}
var file_flyt_legacy_menu_published_proto_depIdxs = []int32{
	1, // 0: flyt.LegacyMenuPublished.restaurant:type_name -> flyt.Ident
	2, // 1: flyt.LegacyMenuPublished.type:type_name -> flyt.MenuType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_flyt_legacy_menu_published_proto_init() }
func file_flyt_legacy_menu_published_proto_init() {
	if File_flyt_legacy_menu_published_proto != nil {
		return
	}
	file_flyt_descriptor_proto_init()
	file_flyt_ident_proto_init()
	file_flyt_legacy_menus_published_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_flyt_legacy_menu_published_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LegacyMenuPublished); i {
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
			RawDescriptor: file_flyt_legacy_menu_published_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flyt_legacy_menu_published_proto_goTypes,
		DependencyIndexes: file_flyt_legacy_menu_published_proto_depIdxs,
		MessageInfos:      file_flyt_legacy_menu_published_proto_msgTypes,
	}.Build()
	File_flyt_legacy_menu_published_proto = out.File
	file_flyt_legacy_menu_published_proto_rawDesc = nil
	file_flyt_legacy_menu_published_proto_goTypes = nil
	file_flyt_legacy_menu_published_proto_depIdxs = nil
}