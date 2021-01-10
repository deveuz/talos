// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.14.0
// source: storage/storage.proto

package storage

import (
	context "context"
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	common "github.com/talos-systems/talos/pkg/machinery/api/common"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Disk represents a disk.
type Disk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Size indicates the disk size in bytes.
	Size uint64 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	// Model idicates the disk model.
	Model string `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	// DeviceName indicates the disk name (e.g. `sda`).
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
}

func (x *Disk) Reset() {
	*x = Disk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Disk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Disk) ProtoMessage() {}

func (x *Disk) ProtoReflect() protoreflect.Message {
	mi := &file_storage_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Disk.ProtoReflect.Descriptor instead.
func (*Disk) Descriptor() ([]byte, []int) {
	return file_storage_storage_proto_rawDescGZIP(), []int{0}
}

func (x *Disk) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Disk) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Disk) GetDeviceName() string {
	if x != nil {
		return x.DeviceName
	}
	return ""
}

// DisksResponse represents the response of the `Disks` RPC.
type DisksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *common.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Disks    []*Disk          `protobuf:"bytes,2,rep,name=disks,proto3" json:"disks,omitempty"`
}

func (x *DisksResponse) Reset() {
	*x = DisksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisksResponse) ProtoMessage() {}

func (x *DisksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_storage_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisksResponse.ProtoReflect.Descriptor instead.
func (*DisksResponse) Descriptor() ([]byte, []int) {
	return file_storage_storage_proto_rawDescGZIP(), []int{1}
}

func (x *DisksResponse) GetMetadata() *common.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *DisksResponse) GetDisks() []*Disk {
	if x != nil {
		return x.Disks
	}
	return nil
}

var File_storage_storage_proto protoreflect.FileDescriptor

var file_storage_storage_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x51, 0x0a, 0x04, 0x44, 0x69, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x62, 0x0a, 0x0d, 0x44, 0x69, 0x73, 0x6b, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x23, 0x0a, 0x05, 0x64, 0x69, 0x73, 0x6b, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x69,
	0x73, 0x6b, 0x52, 0x05, 0x64, 0x69, 0x73, 0x6b, 0x73, 0x32, 0x49, 0x0a, 0x0e, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x05, 0x44,
	0x69, 0x73, 0x6b, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x69, 0x73, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x59, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x42, 0x0a, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x41, 0x70, 0x69, 0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2d, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f,
	0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e,
	0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storage_storage_proto_rawDescOnce sync.Once
	file_storage_storage_proto_rawDescData = file_storage_storage_proto_rawDesc
)

func file_storage_storage_proto_rawDescGZIP() []byte {
	file_storage_storage_proto_rawDescOnce.Do(func() {
		file_storage_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_storage_storage_proto_rawDescData)
	})
	return file_storage_storage_proto_rawDescData
}

var (
	file_storage_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
	file_storage_storage_proto_goTypes  = []interface{}{
		(*Disk)(nil),            // 0: storage.Disk
		(*DisksResponse)(nil),   // 1: storage.DisksResponse
		(*common.Metadata)(nil), // 2: common.Metadata
		(*emptypb.Empty)(nil),   // 3: google.protobuf.Empty
	}
)

var file_storage_storage_proto_depIdxs = []int32{
	2, // 0: storage.DisksResponse.metadata:type_name -> common.Metadata
	0, // 1: storage.DisksResponse.disks:type_name -> storage.Disk
	3, // 2: storage.StorageService.Disks:input_type -> google.protobuf.Empty
	1, // 3: storage.StorageService.Disks:output_type -> storage.DisksResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_storage_storage_proto_init() }
func file_storage_storage_proto_init() {
	if File_storage_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_storage_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Disk); i {
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
		file_storage_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisksResponse); i {
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
			RawDescriptor: file_storage_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storage_storage_proto_goTypes,
		DependencyIndexes: file_storage_storage_proto_depIdxs,
		MessageInfos:      file_storage_storage_proto_msgTypes,
	}.Build()
	File_storage_storage_proto = out.File
	file_storage_storage_proto_rawDesc = nil
	file_storage_storage_proto_goTypes = nil
	file_storage_storage_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StorageServiceClient is the client API for StorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StorageServiceClient interface {
	Disks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DisksResponse, error)
}

type storageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageServiceClient(cc grpc.ClientConnInterface) StorageServiceClient {
	return &storageServiceClient{cc}
}

func (c *storageServiceClient) Disks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DisksResponse, error) {
	out := new(DisksResponse)
	err := c.cc.Invoke(ctx, "/storage.StorageService/Disks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServiceServer is the server API for StorageService service.
type StorageServiceServer interface {
	Disks(context.Context, *emptypb.Empty) (*DisksResponse, error)
}

// UnimplementedStorageServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStorageServiceServer struct {
}

func (*UnimplementedStorageServiceServer) Disks(context.Context, *emptypb.Empty) (*DisksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disks not implemented")
}

func RegisterStorageServiceServer(s *grpc.Server, srv StorageServiceServer) {
	s.RegisterService(&_StorageService_serviceDesc, srv)
}

func _StorageService_Disks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).Disks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.StorageService/Disks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).Disks(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _StorageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "storage.StorageService",
	HandlerType: (*StorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Disks",
			Handler:    _StorageService_Disks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage/storage.proto",
}
