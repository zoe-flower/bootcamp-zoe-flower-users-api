// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: flyt/response-recorded.proto

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

type ResponseRecorded struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The request ID.
	RequestId string `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// The request that was made to us.
	Request *Request `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	// The response that we gave.
	Response *Response `protobuf:"bytes,3,opt,name=response,proto3" json:"response,omitempty"`
	// The capability for which the response was recorded.
	Capability string `protobuf:"bytes,4,opt,name=capability,proto3" json:"capability,omitempty"`
}

func (x *ResponseRecorded) Reset() {
	*x = ResponseRecorded{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flyt_response_recorded_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseRecorded) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseRecorded) ProtoMessage() {}

func (x *ResponseRecorded) ProtoReflect() protoreflect.Message {
	mi := &file_flyt_response_recorded_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseRecorded.ProtoReflect.Descriptor instead.
func (*ResponseRecorded) Descriptor() ([]byte, []int) {
	return file_flyt_response_recorded_proto_rawDescGZIP(), []int{0}
}

func (x *ResponseRecorded) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *ResponseRecorded) GetRequest() *Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *ResponseRecorded) GetResponse() *Response {
	if x != nil {
		return x.Response
	}
	return nil
}

func (x *ResponseRecorded) GetCapability() string {
	if x != nil {
		return x.Capability
	}
	return ""
}

var File_flyt_response_recorded_proto protoreflect.FileDescriptor

var file_flyt_response_recorded_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2d,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04,
	0x66, 0x6c, 0x79, 0x74, 0x1a, 0x15, 0x66, 0x6c, 0x79, 0x74, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x66, 0x6c, 0x79,
	0x74, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x2d, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbd, 0x01, 0x0a, 0x10, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x27, 0x0a,
	0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x6c, 0x79, 0x74, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x70, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x3a, 0x15, 0x82, 0xb5, 0x18, 0x11, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x42, 0x74, 0x0a, 0x08, 0x63, 0x6f, 0x6d,
	0x2e, 0x66, 0x6c, 0x79, 0x74, 0x42, 0x15, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x21,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6c, 0x79, 0x70, 0x61,
	0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x66, 0x6c, 0x79,
	0x74, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0xca, 0x02,
	0x04, 0x46, 0x6c, 0x79, 0x74, 0xe2, 0x02, 0x10, 0x46, 0x6c, 0x79, 0x74, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x04, 0x46, 0x6c, 0x79, 0x74, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flyt_response_recorded_proto_rawDescOnce sync.Once
	file_flyt_response_recorded_proto_rawDescData = file_flyt_response_recorded_proto_rawDesc
)

func file_flyt_response_recorded_proto_rawDescGZIP() []byte {
	file_flyt_response_recorded_proto_rawDescOnce.Do(func() {
		file_flyt_response_recorded_proto_rawDescData = protoimpl.X.CompressGZIP(file_flyt_response_recorded_proto_rawDescData)
	})
	return file_flyt_response_recorded_proto_rawDescData
}

var file_flyt_response_recorded_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_flyt_response_recorded_proto_goTypes = []interface{}{
	(*ResponseRecorded)(nil), // 0: flyt.ResponseRecorded
	(*Request)(nil),          // 1: flyt.Request
	(*Response)(nil),         // 2: flyt.Response
}
var file_flyt_response_recorded_proto_depIdxs = []int32{
	1, // 0: flyt.ResponseRecorded.request:type_name -> flyt.Request
	2, // 1: flyt.ResponseRecorded.response:type_name -> flyt.Response
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_flyt_response_recorded_proto_init() }
func file_flyt_response_recorded_proto_init() {
	if File_flyt_response_recorded_proto != nil {
		return
	}
	file_flyt_descriptor_proto_init()
	file_flyt_callback_recorded_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_flyt_response_recorded_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseRecorded); i {
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
			RawDescriptor: file_flyt_response_recorded_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_flyt_response_recorded_proto_goTypes,
		DependencyIndexes: file_flyt_response_recorded_proto_depIdxs,
		MessageInfos:      file_flyt_response_recorded_proto_msgTypes,
	}.Build()
	File_flyt_response_recorded_proto = out.File
	file_flyt_response_recorded_proto_rawDesc = nil
	file_flyt_response_recorded_proto_goTypes = nil
	file_flyt_response_recorded_proto_depIdxs = nil
}