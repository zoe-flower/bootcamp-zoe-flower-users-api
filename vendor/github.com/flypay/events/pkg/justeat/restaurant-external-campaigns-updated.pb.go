// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: justeat/restaurant-external-campaigns-updated.proto

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

type ExternalCampaign struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExternalCampaignId string `protobuf:"bytes,1,opt,name=external_campaign_id,json=externalCampaignId,proto3" json:"external_campaign_id,omitempty"`
}

func (x *ExternalCampaign) Reset() {
	*x = ExternalCampaign{}
	if protoimpl.UnsafeEnabled {
		mi := &file_justeat_restaurant_external_campaigns_updated_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExternalCampaign) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExternalCampaign) ProtoMessage() {}

func (x *ExternalCampaign) ProtoReflect() protoreflect.Message {
	mi := &file_justeat_restaurant_external_campaigns_updated_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExternalCampaign.ProtoReflect.Descriptor instead.
func (*ExternalCampaign) Descriptor() ([]byte, []int) {
	return file_justeat_restaurant_external_campaigns_updated_proto_rawDescGZIP(), []int{0}
}

func (x *ExternalCampaign) GetExternalCampaignId() string {
	if x != nil {
		return x.ExternalCampaignId
	}
	return ""
}

type RestaurantExternalCampaignsUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantId string `protobuf:"bytes,1,opt,name=restaurant_id,json=RestaurantId,proto3" json:"restaurant_id,omitempty"`
	Tenant       string `protobuf:"bytes,2,opt,name=tenant,json=Tenant,proto3" json:"tenant,omitempty"`
	// e.g 2021-05-05T13:01:49.4976625Z
	Timestamp         string              `protobuf:"bytes,3,opt,name=timestamp,json=Timestamp,proto3" json:"timestamp,omitempty"`
	ExternalCampaigns []*ExternalCampaign `protobuf:"bytes,4,rep,name=external_campaigns,json=ExternalCampaigns,proto3" json:"external_campaigns,omitempty"`
	//service name
	RaisingComponent string `protobuf:"bytes,5,opt,name=raising_component,json=RaisingComponent,proto3" json:"raising_component,omitempty"`
	//Flyt request id
	Conversation string `protobuf:"bytes,6,opt,name=conversation,json=Conversation,proto3" json:"conversation,omitempty"`
}

func (x *RestaurantExternalCampaignsUpdated) Reset() {
	*x = RestaurantExternalCampaignsUpdated{}
	if protoimpl.UnsafeEnabled {
		mi := &file_justeat_restaurant_external_campaigns_updated_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestaurantExternalCampaignsUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestaurantExternalCampaignsUpdated) ProtoMessage() {}

func (x *RestaurantExternalCampaignsUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_justeat_restaurant_external_campaigns_updated_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestaurantExternalCampaignsUpdated.ProtoReflect.Descriptor instead.
func (*RestaurantExternalCampaignsUpdated) Descriptor() ([]byte, []int) {
	return file_justeat_restaurant_external_campaigns_updated_proto_rawDescGZIP(), []int{1}
}

func (x *RestaurantExternalCampaignsUpdated) GetRestaurantId() string {
	if x != nil {
		return x.RestaurantId
	}
	return ""
}

func (x *RestaurantExternalCampaignsUpdated) GetTenant() string {
	if x != nil {
		return x.Tenant
	}
	return ""
}

func (x *RestaurantExternalCampaignsUpdated) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *RestaurantExternalCampaignsUpdated) GetExternalCampaigns() []*ExternalCampaign {
	if x != nil {
		return x.ExternalCampaigns
	}
	return nil
}

func (x *RestaurantExternalCampaignsUpdated) GetRaisingComponent() string {
	if x != nil {
		return x.RaisingComponent
	}
	return ""
}

func (x *RestaurantExternalCampaignsUpdated) GetConversation() string {
	if x != nil {
		return x.Conversation
	}
	return ""
}

var File_justeat_restaurant_external_campaigns_updated_proto protoreflect.FileDescriptor

var file_justeat_restaurant_external_campaigns_updated_proto_rawDesc = []byte{
	0x0a, 0x33, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2f, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75,
	0x72, 0x61, 0x6e, 0x74, 0x2d, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2d, 0x63, 0x61,
	0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x73, 0x2d, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x1a, 0x15,
	0x66, 0x6c, 0x79, 0x74, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x44, 0x0a, 0x10,
	0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e,
	0x12, 0x30, 0x0a, 0x14, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x63, 0x61, 0x6d,
	0x70, 0x61, 0x69, 0x67, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12,
	0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e,
	0x49, 0x64, 0x22, 0xf3, 0x02, 0x0a, 0x22, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e,
	0x74, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67,
	0x6e, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x48, 0x0a, 0x12, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x5f, 0x63, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x52, 0x11, 0x45, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x73, 0x12, 0x2b,
	0x0a, 0x11, 0x72, 0x61, 0x69, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x52, 0x61, 0x69, 0x73, 0x69,
	0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x63,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a,
	0x57, 0x82, 0xb5, 0x18, 0x2d, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x2e, 0x72, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2e, 0x63, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e, 0x73, 0x2e, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0xa2, 0xbb, 0x18, 0x22, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74,
	0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70, 0x61, 0x69, 0x67, 0x6e,
	0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x98, 0x01, 0x0a, 0x0b, 0x63, 0x6f, 0x6d,
	0x2e, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x42, 0x27, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75,
	0x72, 0x61, 0x6e, 0x74, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x43, 0x61, 0x6d, 0x70,
	0x61, 0x69, 0x67, 0x6e, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x6c, 0x79, 0x70, 0x61, 0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0xa2, 0x02, 0x03, 0x4a, 0x58, 0x58, 0xaa,
	0x02, 0x07, 0x4a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0xca, 0x02, 0x07, 0x4a, 0x75, 0x73, 0x74,
	0x65, 0x61, 0x74, 0xe2, 0x02, 0x13, 0x4a, 0x75, 0x73, 0x74, 0x65, 0x61, 0x74, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x07, 0x4a, 0x75, 0x73, 0x74,
	0x65, 0x61, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_justeat_restaurant_external_campaigns_updated_proto_rawDescOnce sync.Once
	file_justeat_restaurant_external_campaigns_updated_proto_rawDescData = file_justeat_restaurant_external_campaigns_updated_proto_rawDesc
)

func file_justeat_restaurant_external_campaigns_updated_proto_rawDescGZIP() []byte {
	file_justeat_restaurant_external_campaigns_updated_proto_rawDescOnce.Do(func() {
		file_justeat_restaurant_external_campaigns_updated_proto_rawDescData = protoimpl.X.CompressGZIP(file_justeat_restaurant_external_campaigns_updated_proto_rawDescData)
	})
	return file_justeat_restaurant_external_campaigns_updated_proto_rawDescData
}

var file_justeat_restaurant_external_campaigns_updated_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_justeat_restaurant_external_campaigns_updated_proto_goTypes = []interface{}{
	(*ExternalCampaign)(nil),                   // 0: justeat.ExternalCampaign
	(*RestaurantExternalCampaignsUpdated)(nil), // 1: justeat.RestaurantExternalCampaignsUpdated
}
var file_justeat_restaurant_external_campaigns_updated_proto_depIdxs = []int32{
	0, // 0: justeat.RestaurantExternalCampaignsUpdated.external_campaigns:type_name -> justeat.ExternalCampaign
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_justeat_restaurant_external_campaigns_updated_proto_init() }
func file_justeat_restaurant_external_campaigns_updated_proto_init() {
	if File_justeat_restaurant_external_campaigns_updated_proto != nil {
		return
	}
	file_justeat_options_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_justeat_restaurant_external_campaigns_updated_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExternalCampaign); i {
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
		file_justeat_restaurant_external_campaigns_updated_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestaurantExternalCampaignsUpdated); i {
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
			RawDescriptor: file_justeat_restaurant_external_campaigns_updated_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_justeat_restaurant_external_campaigns_updated_proto_goTypes,
		DependencyIndexes: file_justeat_restaurant_external_campaigns_updated_proto_depIdxs,
		MessageInfos:      file_justeat_restaurant_external_campaigns_updated_proto_msgTypes,
	}.Build()
	File_justeat_restaurant_external_campaigns_updated_proto = out.File
	file_justeat_restaurant_external_campaigns_updated_proto_rawDesc = nil
	file_justeat_restaurant_external_campaigns_updated_proto_goTypes = nil
	file_justeat_restaurant_external_campaigns_updated_proto_depIdxs = nil
}
