// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: flyt/restaurant-ordering-capacity-updated.proto

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

type RestaurantOrderingCapacityUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Define your proto fields here
	Restaurant              *RestaurantIdent `protobuf:"bytes,1,opt,name=restaurant,proto3" json:"restaurant,omitempty"`
	CookTimeEarliestMinutes int32            `protobuf:"varint,2,opt,name=cook_time_earliest_minutes,json=cookTimeEarliestMinutes,proto3" json:"cook_time_earliest_minutes,omitempty"`
	CookTimeIntervalMinutes int32            `protobuf:"varint,3,opt,name=cook_time_interval_minutes,json=cookTimeIntervalMinutes,proto3" json:"cook_time_interval_minutes,omitempty"`
}

func (x *RestaurantOrderingCapacityUpdated) Reset() {
	*x = RestaurantOrderingCapacityUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_restaurant_ordering_capacity_updated_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestaurantOrderingCapacityUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestaurantOrderingCapacityUpdated) ProtoMessage() {}

func (x *RestaurantOrderingCapacityUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_restaurant_ordering_capacity_updated_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestaurantOrderingCapacityUpdated.ProtoReflect.Descriptor instead.
func (*RestaurantOrderingCapacityUpdated) Descriptor() ([]byte, []int) {
	return file_flyt_restaurant_ordering_capacity_updated_proto_rawDescGZIP(), []int{0}
}

func (x *RestaurantOrderingCapacityUpdated) GetRestaurant() *RestaurantIdent {
	if x != nil {
		return x.Restaurant
	}
	return nil
}

func (x *RestaurantOrderingCapacityUpdated) GetCookTimeEarliestMinutes() int32 {
	if x != nil {
		return x.CookTimeEarliestMinutes
	}
	return 0
}

func (x *RestaurantOrderingCapacityUpdated) GetCookTimeIntervalMinutes() int32 {
	if x != nil {
		return x.CookTimeIntervalMinutes
	}
	return 0
}

var File_flyt_restaurant_ordering_capacity_updated_proto protoreflect.FileDescriptor

var file_flyt_restaurant_ordering_capacity_updated_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e,
	0x74, 0x2d, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2d, 0x63, 0x61, 0x70, 0x61, 0x63,
	0x69, 0x74, 0x79, 0x2d, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x66, 0x6c, 0x79, 0x74, 0x1a, 0x15, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x66, 0x6c, 0x79, 0x74, 0x2f, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfe, 0x01, 0x0a, 0x21, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75,
	0x72, 0x61, 0x6e, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x70, 0x61,
	0x63, 0x69, 0x74, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x35, 0x0a, 0x0a, 0x72,
	0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e,
	0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61,
	0x6e, 0x74, 0x12, 0x3b, 0x0a, 0x1a, 0x63, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x65, 0x61, 0x72, 0x6c, 0x69, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x17, 0x63, 0x6f, 0x6f, 0x6b, 0x54, 0x69, 0x6d, 0x65,
	0x45, 0x61, 0x72, 0x6c, 0x69, 0x65, 0x73, 0x74, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x12,
	0x3b, 0x0a, 0x1a, 0x63, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x5f, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x17, 0x63, 0x6f, 0x6f, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x3a, 0x28, 0x82, 0xb5,
	0x18, 0x24, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x2e, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x85, 0x01, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x2e, 0x66,
	0x6c, 0x79, 0x74, 0x42, 0x26, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x21, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x79, 0x70, 0x61, 0x79,
	0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x66, 0x6c, 0x79, 0x74,
	0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0xca, 0x02, 0x04,
	0x46, 0x6c, 0x79, 0x74, 0xe2, 0x02, 0x10, 0x46, 0x6c, 0x79, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flyt_restaurant_ordering_capacity_updated_proto_rawDescOnce sync.Once
	file_flyt_restaurant_ordering_capacity_updated_proto_rawDescData = file_flyt_restaurant_ordering_capacity_updated_proto_rawDesc
)

func file_flyt_restaurant_ordering_capacity_updated_proto_rawDescGZIP() []byte {
	file_flyt_restaurant_ordering_capacity_updated_proto_rawDescOnce.Do(func() {
		file_flyt_restaurant_ordering_capacity_updated_proto_rawDescData = protoimpl.X.CompressGZIP(file_flyt_restaurant_ordering_capacity_updated_proto_rawDescData)
	})
	return file_flyt_restaurant_ordering_capacity_updated_proto_rawDescData
}

var file_flyt_restaurant_ordering_capacity_updated_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_flyt_restaurant_ordering_capacity_updated_proto_goTypes = []interface{}{
	(*RestaurantOrderingCapacityUpdated)(nil), // 0: flyt.RestaurantOrderingCapacityUpdated
	(*RestaurantIdent)(nil),                   // 1: flyt.RestaurantIdent
}
var file_flyt_restaurant_ordering_capacity_updated_proto_depIdxs = []int32{
	1, // 0: flyt.RestaurantOrderingCapacityUpdated.restaurant:type_name -> flyt.RestaurantIdent
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_flyt_restaurant_ordering_capacity_updated_proto_init() }
func file_flyt_restaurant_ordering_capacity_updated_proto_init() {
	if File_flyt_restaurant_ordering_capacity_updated_proto != nil {
		return
	}
	file_flyt_descriptor_proto_init()
	file_flyt_restaurant_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_flyt_restaurant_ordering_capacity_updated_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestaurantOrderingCapacityUpdated); i {
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
			RawDescriptor: file_flyt_restaurant_ordering_capacity_updated_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flyt_restaurant_ordering_capacity_updated_proto_goTypes,
		DependencyIndexes: file_flyt_restaurant_ordering_capacity_updated_proto_depIdxs,
		MessageInfos:      file_flyt_restaurant_ordering_capacity_updated_proto_msgTypes,
	}.Build()
	File_flyt_restaurant_ordering_capacity_updated_proto = out.File
	file_flyt_restaurant_ordering_capacity_updated_proto_rawDesc = nil
	file_flyt_restaurant_ordering_capacity_updated_proto_goTypes = nil
	file_flyt_restaurant_ordering_capacity_updated_proto_depIdxs = nil
}
