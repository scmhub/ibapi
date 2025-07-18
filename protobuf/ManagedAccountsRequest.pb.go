// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: ManagedAccountsRequest.proto

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

type ManagedAccountsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ManagedAccountsRequest) Reset() {
	*x = ManagedAccountsRequest{}
	mi := &file_ManagedAccountsRequest_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ManagedAccountsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManagedAccountsRequest) ProtoMessage() {}

func (x *ManagedAccountsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ManagedAccountsRequest_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManagedAccountsRequest.ProtoReflect.Descriptor instead.
func (*ManagedAccountsRequest) Descriptor() ([]byte, []int) {
	return file_ManagedAccountsRequest_proto_rawDescGZIP(), []int{0}
}

var File_ManagedAccountsRequest_proto protoreflect.FileDescriptor

var file_ManagedAccountsRequest_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x18, 0x0a, 0x16, 0x4d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_ManagedAccountsRequest_proto_rawDescOnce sync.Once
	file_ManagedAccountsRequest_proto_rawDescData []byte
)

func file_ManagedAccountsRequest_proto_rawDescGZIP() []byte {
	file_ManagedAccountsRequest_proto_rawDescOnce.Do(func() {
		file_ManagedAccountsRequest_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_ManagedAccountsRequest_proto_rawDesc), len(file_ManagedAccountsRequest_proto_rawDesc)))
	})
	return file_ManagedAccountsRequest_proto_rawDescData
}

var file_ManagedAccountsRequest_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ManagedAccountsRequest_proto_goTypes = []any{
	(*ManagedAccountsRequest)(nil), // 0: protobuf.ManagedAccountsRequest
}
var file_ManagedAccountsRequest_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ManagedAccountsRequest_proto_init() }
func file_ManagedAccountsRequest_proto_init() {
	if File_ManagedAccountsRequest_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_ManagedAccountsRequest_proto_rawDesc), len(file_ManagedAccountsRequest_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ManagedAccountsRequest_proto_goTypes,
		DependencyIndexes: file_ManagedAccountsRequest_proto_depIdxs,
		MessageInfos:      file_ManagedAccountsRequest_proto_msgTypes,
	}.Build()
	File_ManagedAccountsRequest_proto = out.File
	file_ManagedAccountsRequest_proto_goTypes = nil
	file_ManagedAccountsRequest_proto_depIdxs = nil
}
