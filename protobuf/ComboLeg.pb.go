// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: ComboLeg.proto

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

type ComboLeg struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	ConId              *int32                 `protobuf:"varint,1,opt,name=conId,proto3,oneof" json:"conId,omitempty"`
	Ratio              *int32                 `protobuf:"varint,2,opt,name=ratio,proto3,oneof" json:"ratio,omitempty"`
	Action             *string                `protobuf:"bytes,3,opt,name=action,proto3,oneof" json:"action,omitempty"`
	Exchange           *string                `protobuf:"bytes,4,opt,name=exchange,proto3,oneof" json:"exchange,omitempty"`
	OpenClose          *int32                 `protobuf:"varint,5,opt,name=openClose,proto3,oneof" json:"openClose,omitempty"`
	ShortSalesSlot     *int32                 `protobuf:"varint,6,opt,name=shortSalesSlot,proto3,oneof" json:"shortSalesSlot,omitempty"`
	DesignatedLocation *string                `protobuf:"bytes,7,opt,name=designatedLocation,proto3,oneof" json:"designatedLocation,omitempty"`
	ExemptCode         *int32                 `protobuf:"varint,8,opt,name=exemptCode,proto3,oneof" json:"exemptCode,omitempty"`
	PerLegPrice        *float64               `protobuf:"fixed64,9,opt,name=perLegPrice,proto3,oneof" json:"perLegPrice,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *ComboLeg) Reset() {
	*x = ComboLeg{}
	mi := &file_ComboLeg_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ComboLeg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComboLeg) ProtoMessage() {}

func (x *ComboLeg) ProtoReflect() protoreflect.Message {
	mi := &file_ComboLeg_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComboLeg.ProtoReflect.Descriptor instead.
func (*ComboLeg) Descriptor() ([]byte, []int) {
	return file_ComboLeg_proto_rawDescGZIP(), []int{0}
}

func (x *ComboLeg) GetConId() int32 {
	if x != nil && x.ConId != nil {
		return *x.ConId
	}
	return 0
}

func (x *ComboLeg) GetRatio() int32 {
	if x != nil && x.Ratio != nil {
		return *x.Ratio
	}
	return 0
}

func (x *ComboLeg) GetAction() string {
	if x != nil && x.Action != nil {
		return *x.Action
	}
	return ""
}

func (x *ComboLeg) GetExchange() string {
	if x != nil && x.Exchange != nil {
		return *x.Exchange
	}
	return ""
}

func (x *ComboLeg) GetOpenClose() int32 {
	if x != nil && x.OpenClose != nil {
		return *x.OpenClose
	}
	return 0
}

func (x *ComboLeg) GetShortSalesSlot() int32 {
	if x != nil && x.ShortSalesSlot != nil {
		return *x.ShortSalesSlot
	}
	return 0
}

func (x *ComboLeg) GetDesignatedLocation() string {
	if x != nil && x.DesignatedLocation != nil {
		return *x.DesignatedLocation
	}
	return ""
}

func (x *ComboLeg) GetExemptCode() int32 {
	if x != nil && x.ExemptCode != nil {
		return *x.ExemptCode
	}
	return 0
}

func (x *ComboLeg) GetPerLegPrice() float64 {
	if x != nil && x.PerLegPrice != nil {
		return *x.PerLegPrice
	}
	return 0
}

var File_ComboLeg_proto protoreflect.FileDescriptor

var file_ComboLeg_proto_rawDesc = string([]byte{
	0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x62, 0x6f, 0x4c, 0x65, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0xd2, 0x03, 0x0a, 0x08, 0x43,
	0x6f, 0x6d, 0x62, 0x6f, 0x4c, 0x65, 0x67, 0x12, 0x19, 0x0a, 0x05, 0x63, 0x6f, 0x6e, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x05, 0x63, 0x6f, 0x6e, 0x49, 0x64, 0x88,
	0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x01, 0x52, 0x05, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x65, 0x78,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x08,
	0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6f,
	0x70, 0x65, 0x6e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x48, 0x04,
	0x52, 0x09, 0x6f, 0x70, 0x65, 0x6e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x88, 0x01, 0x01, 0x12, 0x2b,
	0x0a, 0x0e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x53, 0x61, 0x6c, 0x65, 0x73, 0x53, 0x6c, 0x6f, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x48, 0x05, 0x52, 0x0e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x53,
	0x61, 0x6c, 0x65, 0x73, 0x53, 0x6c, 0x6f, 0x74, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x12, 0x64,
	0x65, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52, 0x12, 0x64, 0x65, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x65, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x12, 0x23, 0x0a, 0x0a, 0x65, 0x78, 0x65, 0x6d, 0x70, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x07, 0x52, 0x0a, 0x65, 0x78, 0x65, 0x6d, 0x70, 0x74, 0x43, 0x6f,
	0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x4c, 0x65, 0x67, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x48, 0x08, 0x52, 0x0b, 0x70, 0x65,
	0x72, 0x4c, 0x65, 0x67, 0x50, 0x72, 0x69, 0x63, 0x65, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06,
	0x5f, 0x63, 0x6f, 0x6e, 0x49, 0x64, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f,
	0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6f, 0x70, 0x65,
	0x6e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x53, 0x61, 0x6c, 0x65, 0x73, 0x53, 0x6c, 0x6f, 0x74, 0x42, 0x15, 0x0a, 0x13, 0x5f, 0x64, 0x65,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x65, 0x64, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x65, 0x78, 0x65, 0x6d, 0x70, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x42,
	0x0e, 0x0a, 0x0c, 0x5f, 0x70, 0x65, 0x72, 0x4c, 0x65, 0x67, 0x50, 0x72, 0x69, 0x63, 0x65, 0x42,
	0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_ComboLeg_proto_rawDescOnce sync.Once
	file_ComboLeg_proto_rawDescData []byte
)

func file_ComboLeg_proto_rawDescGZIP() []byte {
	file_ComboLeg_proto_rawDescOnce.Do(func() {
		file_ComboLeg_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_ComboLeg_proto_rawDesc), len(file_ComboLeg_proto_rawDesc)))
	})
	return file_ComboLeg_proto_rawDescData
}

var file_ComboLeg_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ComboLeg_proto_goTypes = []any{
	(*ComboLeg)(nil), // 0: protobuf.ComboLeg
}
var file_ComboLeg_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ComboLeg_proto_init() }
func file_ComboLeg_proto_init() {
	if File_ComboLeg_proto != nil {
		return
	}
	file_ComboLeg_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_ComboLeg_proto_rawDesc), len(file_ComboLeg_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ComboLeg_proto_goTypes,
		DependencyIndexes: file_ComboLeg_proto_depIdxs,
		MessageInfos:      file_ComboLeg_proto_msgTypes,
	}.Build()
	File_ComboLeg_proto = out.File
	file_ComboLeg_proto_goTypes = nil
	file_ComboLeg_proto_depIdxs = nil
}
