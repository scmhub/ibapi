// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: OrderStatus.proto

package protobuf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderId       *int32                 `protobuf:"varint,1,opt,name=orderId,proto3,oneof" json:"orderId,omitempty"`
	Status        *string                `protobuf:"bytes,2,opt,name=status,proto3,oneof" json:"status,omitempty"`
	Filled        *string                `protobuf:"bytes,3,opt,name=filled,proto3,oneof" json:"filled,omitempty"`
	Remaining     *string                `protobuf:"bytes,4,opt,name=remaining,proto3,oneof" json:"remaining,omitempty"`
	AvgFillPrice  *float64               `protobuf:"fixed64,5,opt,name=avgFillPrice,proto3,oneof" json:"avgFillPrice,omitempty"`
	PermId        *int64                 `protobuf:"varint,6,opt,name=permId,proto3,oneof" json:"permId,omitempty"`
	ParentId      *int32                 `protobuf:"varint,7,opt,name=parentId,proto3,oneof" json:"parentId,omitempty"`
	LastFillPrice *float64               `protobuf:"fixed64,8,opt,name=lastFillPrice,proto3,oneof" json:"lastFillPrice,omitempty"`
	ClientId      *int32                 `protobuf:"varint,9,opt,name=clientId,proto3,oneof" json:"clientId,omitempty"`
	WhyHeld       *string                `protobuf:"bytes,10,opt,name=whyHeld,proto3,oneof" json:"whyHeld,omitempty"`
	MktCapPrice   *float64               `protobuf:"fixed64,11,opt,name=mktCapPrice,proto3,oneof" json:"mktCapPrice,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderStatus) Reset() {
	*x = OrderStatus{}
	mi := &file_OrderStatus_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderStatus) ProtoMessage() {}

func (x *OrderStatus) ProtoReflect() protoreflect.Message {
	mi := &file_OrderStatus_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderStatus.ProtoReflect.Descriptor instead.
func (*OrderStatus) Descriptor() ([]byte, []int) {
	return file_OrderStatus_proto_rawDescGZIP(), []int{0}
}

func (x *OrderStatus) GetOrderId() int32 {
	if x != nil && x.OrderId != nil {
		return *x.OrderId
	}
	return 0
}

func (x *OrderStatus) GetStatus() string {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return ""
}

func (x *OrderStatus) GetFilled() string {
	if x != nil && x.Filled != nil {
		return *x.Filled
	}
	return ""
}

func (x *OrderStatus) GetRemaining() string {
	if x != nil && x.Remaining != nil {
		return *x.Remaining
	}
	return ""
}

func (x *OrderStatus) GetAvgFillPrice() float64 {
	if x != nil && x.AvgFillPrice != nil {
		return *x.AvgFillPrice
	}
	return 0
}

func (x *OrderStatus) GetPermId() int64 {
	if x != nil && x.PermId != nil {
		return *x.PermId
	}
	return 0
}

func (x *OrderStatus) GetParentId() int32 {
	if x != nil && x.ParentId != nil {
		return *x.ParentId
	}
	return 0
}

func (x *OrderStatus) GetLastFillPrice() float64 {
	if x != nil && x.LastFillPrice != nil {
		return *x.LastFillPrice
	}
	return 0
}

func (x *OrderStatus) GetClientId() int32 {
	if x != nil && x.ClientId != nil {
		return *x.ClientId
	}
	return 0
}

func (x *OrderStatus) GetWhyHeld() string {
	if x != nil && x.WhyHeld != nil {
		return *x.WhyHeld
	}
	return ""
}

func (x *OrderStatus) GetMktCapPrice() float64 {
	if x != nil && x.MktCapPrice != nil {
		return *x.MktCapPrice
	}
	return 0
}

var File_OrderStatus_proto protoreflect.FileDescriptor

var file_OrderStatus_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x96, 0x04,
	0x0a, 0x0b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a,
	0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00,
	0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x66, 0x69, 0x6c,
	0x6c, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x06, 0x66, 0x69, 0x6c,
	0x6c, 0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e,
	0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x09, 0x72, 0x65, 0x6d,
	0x61, 0x69, 0x6e, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0c, 0x61, 0x76, 0x67,
	0x46, 0x69, 0x6c, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x48,
	0x04, 0x52, 0x0c, 0x61, 0x76, 0x67, 0x46, 0x69, 0x6c, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x88,
	0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x6d, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x48, 0x05, 0x52, 0x06, 0x70, 0x65, 0x72, 0x6d, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x1f, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x06, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01,
	0x12, 0x29, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x6c, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x48, 0x07, 0x52, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x46,
	0x69, 0x6c, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x48, 0x08, 0x52,
	0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07,
	0x77, 0x68, 0x79, 0x48, 0x65, 0x6c, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x09, 0x52,
	0x07, 0x77, 0x68, 0x79, 0x48, 0x65, 0x6c, 0x64, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x6d,
	0x6b, 0x74, 0x43, 0x61, 0x70, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01,
	0x48, 0x0a, 0x52, 0x0b, 0x6d, 0x6b, 0x74, 0x43, 0x61, 0x70, 0x50, 0x72, 0x69, 0x63, 0x65, 0x88,
	0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x42, 0x09,
	0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x66, 0x69,
	0x6c, 0x6c, 0x65, 0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x69,
	0x6e, 0x67, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x61, 0x76, 0x67, 0x46, 0x69, 0x6c, 0x6c, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x70, 0x65, 0x72, 0x6d, 0x49, 0x64, 0x42, 0x0b,
	0x0a, 0x09, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x42, 0x10, 0x0a, 0x0e, 0x5f,
	0x6c, 0x61, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x42, 0x0b, 0x0a,
	0x09, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x77,
	0x68, 0x79, 0x48, 0x65, 0x6c, 0x64, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x6d, 0x6b, 0x74, 0x43, 0x61,
	0x70, 0x50, 0x72, 0x69, 0x63, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_OrderStatus_proto_rawDescOnce sync.Once
	file_OrderStatus_proto_rawDescData []byte
)

func file_OrderStatus_proto_rawDescGZIP() []byte {
	file_OrderStatus_proto_rawDescOnce.Do(func() {
		file_OrderStatus_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_OrderStatus_proto_rawDesc), len(file_OrderStatus_proto_rawDesc)))
	})
	return file_OrderStatus_proto_rawDescData
}

var file_OrderStatus_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_OrderStatus_proto_goTypes = []any{
	(*OrderStatus)(nil), // 0: protobuf.OrderStatus
}
var file_OrderStatus_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_OrderStatus_proto_init() }
func file_OrderStatus_proto_init() {
	if File_OrderStatus_proto != nil {
		return
	}
	file_OrderStatus_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_OrderStatus_proto_rawDesc), len(file_OrderStatus_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_OrderStatus_proto_goTypes,
		DependencyIndexes: file_OrderStatus_proto_depIdxs,
		MessageInfos:      file_OrderStatus_proto_msgTypes,
	}.Build()
	File_OrderStatus_proto = out.File
	file_OrderStatus_proto_goTypes = nil
	file_OrderStatus_proto_depIdxs = nil
}
