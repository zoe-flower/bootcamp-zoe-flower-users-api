// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: flyt/order-return-failed.proto

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

// Define your proto fields here
type OrderReturnFailed_ErrorCode int32

const (
	OrderReturnFailed_UNKNOWN            OrderReturnFailed_ErrorCode = 0
	OrderReturnFailed_NOT_SUPPORTED      OrderReturnFailed_ErrorCode = 1
	OrderReturnFailed_MALFORMED_REQUEST  OrderReturnFailed_ErrorCode = 2
	OrderReturnFailed_AUTH_FAILED        OrderReturnFailed_ErrorCode = 3
	OrderReturnFailed_CONNECTIVITY_ISSUE OrderReturnFailed_ErrorCode = 4
	OrderReturnFailed_TIMEOUT            OrderReturnFailed_ErrorCode = 5
)

// Enum value maps for OrderReturnFailed_ErrorCode.
var (
	OrderReturnFailed_ErrorCode_name = map[int32]string{
		0: "UNKNOWN",
		1: "NOT_SUPPORTED",
		2: "MALFORMED_REQUEST",
		3: "AUTH_FAILED",
		4: "CONNECTIVITY_ISSUE",
		5: "TIMEOUT",
	}
	OrderReturnFailed_ErrorCode_value = map[string]int32{
		"UNKNOWN":            0,
		"NOT_SUPPORTED":      1,
		"MALFORMED_REQUEST":  2,
		"AUTH_FAILED":        3,
		"CONNECTIVITY_ISSUE": 4,
		"TIMEOUT":            5,
	}
)

func (x OrderReturnFailed_ErrorCode) Enum() *OrderReturnFailed_ErrorCode {
	p := new(OrderReturnFailed_ErrorCode)
	*p = x
	return p
}

func (x OrderReturnFailed_ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderReturnFailed_ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_flyt_order_return_failed_proto_enumTypes[0].Descriptor()
}

func (OrderReturnFailed_ErrorCode) Type() protoreflect.EnumType {
	return &file_flyt_order_return_failed_proto_enumTypes[0]
}

func (x OrderReturnFailed_ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderReturnFailed_ErrorCode.Descriptor instead.
func (OrderReturnFailed_ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_flyt_order_return_failed_proto_rawDescGZIP(), []int{0, 0}
}

type OrderReturnFailed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId    string                   `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	Error      *OrderReturnFailed_Error `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Restaurant *Ident                   `protobuf:"bytes,3,opt,name=restaurant,proto3" json:"restaurant,omitempty"`
}

func (x *OrderReturnFailed) Reset() {
	*x = OrderReturnFailed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_order_return_failed_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderReturnFailed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderReturnFailed) ProtoMessage() {}

func (x *OrderReturnFailed) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_order_return_failed_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderReturnFailed.ProtoReflect.Descriptor instead.
func (*OrderReturnFailed) Descriptor() ([]byte, []int) {
	return file_flyt_order_return_failed_proto_rawDescGZIP(), []int{0}
}

func (x *OrderReturnFailed) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *OrderReturnFailed) GetError() *OrderReturnFailed_Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *OrderReturnFailed) GetRestaurant() *Ident {
	if x != nil {
		return x.Restaurant
	}
	return nil
}

type OrderReturnFailed_Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    OrderReturnFailed_ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=flyt.OrderReturnFailed_ErrorCode" json:"code,omitempty"` // Define error's type
	Message string                      `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`                                  // Description of the error
}

func (x *OrderReturnFailed_Error) Reset() {
	*x = OrderReturnFailed_Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_order_return_failed_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderReturnFailed_Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderReturnFailed_Error) ProtoMessage() {}

func (x *OrderReturnFailed_Error) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_order_return_failed_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderReturnFailed_Error.ProtoReflect.Descriptor instead.
func (*OrderReturnFailed_Error) Descriptor() ([]byte, []int) {
	return file_flyt_order_return_failed_proto_rawDescGZIP(), []int{0, 0}
}

func (x *OrderReturnFailed_Error) GetCode() OrderReturnFailed_ErrorCode {
	if x != nil {
		return x.Code
	}
	return OrderReturnFailed_UNKNOWN
}

func (x *OrderReturnFailed_Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_flyt_order_return_failed_proto protoreflect.FileDescriptor

var file_flyt_order_return_failed_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2d, 0x72, 0x65, 0x74,
	0x75, 0x72, 0x6e, 0x2d, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x66, 0x6c, 0x79, 0x74, 0x1a, 0x15, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x66,
	0x6c, 0x79, 0x74, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xfd, 0x02, 0x0a, 0x11, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x46,
	0x61, 0x69, 0x6c, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x33, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x2b, 0x0a, 0x0a, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72,
	0x61, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x66, 0x6c, 0x79, 0x74,
	0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61,
	0x6e, 0x74, 0x1a, 0x58, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x35, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x66, 0x6c, 0x79, 0x74,
	0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x46, 0x61, 0x69, 0x6c,
	0x65, 0x64, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x78, 0x0a, 0x09,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x55,
	0x50, 0x50, 0x4f, 0x52, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x4d, 0x41, 0x4c,
	0x46, 0x4f, 0x52, 0x4d, 0x45, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x02,
	0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x55, 0x54, 0x48, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x49, 0x56, 0x49, 0x54,
	0x59, 0x5f, 0x49, 0x53, 0x53, 0x55, 0x45, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x49, 0x4d,
	0x45, 0x4f, 0x55, 0x54, 0x10, 0x05, 0x3a, 0x17, 0x82, 0xb5, 0x18, 0x13, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x2e, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x2e, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x42,
	0x75, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x42, 0x16, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x6c, 0x79, 0x70, 0x61, 0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x66, 0x6c, 0x79, 0x74, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02,
	0x04, 0x46, 0x6c, 0x79, 0x74, 0xca, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0xe2, 0x02, 0x10, 0x46,
	0x6c, 0x79, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flyt_order_return_failed_proto_rawDescOnce sync.Once
	file_flyt_order_return_failed_proto_rawDescData = file_flyt_order_return_failed_proto_rawDesc
)

func file_flyt_order_return_failed_proto_rawDescGZIP() []byte {
	file_flyt_order_return_failed_proto_rawDescOnce.Do(func() {
		file_flyt_order_return_failed_proto_rawDescData = protoimpl.X.CompressGZIP(file_flyt_order_return_failed_proto_rawDescData)
	})
	return file_flyt_order_return_failed_proto_rawDescData
}

var file_flyt_order_return_failed_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_flyt_order_return_failed_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_flyt_order_return_failed_proto_goTypes = []interface{}{
	(OrderReturnFailed_ErrorCode)(0), // 0: flyt.OrderReturnFailed.ErrorCode
	(*OrderReturnFailed)(nil),        // 1: flyt.OrderReturnFailed
	(*OrderReturnFailed_Error)(nil),  // 2: flyt.OrderReturnFailed.Error
	(*Ident)(nil),                    // 3: flyt.Ident
}
var file_flyt_order_return_failed_proto_depIdxs = []int32{
	2, // 0: flyt.OrderReturnFailed.error:type_name -> flyt.OrderReturnFailed.Error
	3, // 1: flyt.OrderReturnFailed.restaurant:type_name -> flyt.Ident
	0, // 2: flyt.OrderReturnFailed.Error.code:type_name -> flyt.OrderReturnFailed.ErrorCode
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_flyt_order_return_failed_proto_init() }
func file_flyt_order_return_failed_proto_init() {
	if File_flyt_order_return_failed_proto != nil {
		return
	}
	file_flyt_descriptor_proto_init()
	file_flyt_ident_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_flyt_order_return_failed_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderReturnFailed); i {
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
		file_flyt_order_return_failed_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OrderReturnFailed_Error); i {
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
			RawDescriptor: file_flyt_order_return_failed_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flyt_order_return_failed_proto_goTypes,
		DependencyIndexes: file_flyt_order_return_failed_proto_depIdxs,
		EnumInfos:         file_flyt_order_return_failed_proto_enumTypes,
		MessageInfos:      file_flyt_order_return_failed_proto_msgTypes,
	}.Build()
	File_flyt_order_return_failed_proto = out.File
	file_flyt_order_return_failed_proto_rawDesc = nil
	file_flyt_order_return_failed_proto_goTypes = nil
	file_flyt_order_return_failed_proto_depIdxs = nil
}