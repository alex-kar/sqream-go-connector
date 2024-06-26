// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: authentication_type.proto

package proto

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

type AuthenticationType int32

const (
	AuthenticationType_AUTHENTICATION_TYPE_UNSPECIFIED AuthenticationType = 0
	AuthenticationType_AUTHENTICATION_TYPE_INTERNAL    AuthenticationType = 1
	AuthenticationType_AUTHENTICATION_TYPE_IDP         AuthenticationType = 2
)

// Enum value maps for AuthenticationType.
var (
	AuthenticationType_name = map[int32]string{
		0: "AUTHENTICATION_TYPE_UNSPECIFIED",
		1: "AUTHENTICATION_TYPE_INTERNAL",
		2: "AUTHENTICATION_TYPE_IDP",
	}
	AuthenticationType_value = map[string]int32{
		"AUTHENTICATION_TYPE_UNSPECIFIED": 0,
		"AUTHENTICATION_TYPE_INTERNAL":    1,
		"AUTHENTICATION_TYPE_IDP":         2,
	}
)

func (x AuthenticationType) Enum() *AuthenticationType {
	p := new(AuthenticationType)
	*p = x
	return p
}

func (x AuthenticationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthenticationType) Descriptor() protoreflect.EnumDescriptor {
	return file_authentication_type_proto_enumTypes[0].Descriptor()
}

func (AuthenticationType) Type() protoreflect.EnumType {
	return &file_authentication_type_proto_enumTypes[0]
}

func (x AuthenticationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthenticationType.Descriptor instead.
func (AuthenticationType) EnumDescriptor() ([]byte, []int) {
	return file_authentication_type_proto_rawDescGZIP(), []int{0}
}

var File_authentication_type_proto protoreflect.FileDescriptor

var file_authentication_type_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x71, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2a, 0x78, 0x0a, 0x12, 0x41, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x23, 0x0a, 0x1f, 0x41, 0x55, 0x54, 0x48, 0x45, 0x4e, 0x54, 0x49, 0x43, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x20, 0x0a, 0x1c, 0x41, 0x55, 0x54, 0x48, 0x45, 0x4e, 0x54,
	0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54,
	0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x41, 0x55, 0x54, 0x48, 0x45,
	0x4e, 0x54, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49,
	0x44, 0x50, 0x10, 0x02, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_authentication_type_proto_rawDescOnce sync.Once
	file_authentication_type_proto_rawDescData = file_authentication_type_proto_rawDesc
)

func file_authentication_type_proto_rawDescGZIP() []byte {
	file_authentication_type_proto_rawDescOnce.Do(func() {
		file_authentication_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_authentication_type_proto_rawDescData)
	})
	return file_authentication_type_proto_rawDescData
}

var file_authentication_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_authentication_type_proto_goTypes = []interface{}{
	(AuthenticationType)(0), // 0: com.sqream.cloud.generated.v1.AuthenticationType
}
var file_authentication_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_authentication_type_proto_init() }
func file_authentication_type_proto_init() {
	if File_authentication_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_authentication_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_authentication_type_proto_goTypes,
		DependencyIndexes: file_authentication_type_proto_depIdxs,
		EnumInfos:         file_authentication_type_proto_enumTypes,
	}.Build()
	File_authentication_type_proto = out.File
	file_authentication_type_proto_rawDesc = nil
	file_authentication_type_proto_goTypes = nil
	file_authentication_type_proto_depIdxs = nil
}
