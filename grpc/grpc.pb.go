// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: grpc.proto

package grpc

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

type Timer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Project string   `protobuf:"bytes,2,opt,name=project,proto3" json:"project,omitempty"`
	Task    string   `protobuf:"bytes,3,opt,name=task,proto3" json:"task,omitempty"`
	Tags    []string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	Start   string   `protobuf:"bytes,5,opt,name=start,proto3" json:"start,omitempty"`
	Stop    string   `protobuf:"bytes,6,opt,name=stop,proto3" json:"stop,omitempty"`
}

func (x *Timer) Reset() {
	*x = Timer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Timer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timer) ProtoMessage() {}

func (x *Timer) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timer.ProtoReflect.Descriptor instead.
func (*Timer) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *Timer) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Timer) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *Timer) GetTask() string {
	if x != nil {
		return x.Task
	}
	return ""
}

func (x *Timer) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Timer) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *Timer) GetStop() string {
	if x != nil {
		return x.Stop
	}
	return ""
}

type StartParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project   string   `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Task      string   `protobuf:"bytes,2,opt,name=task,proto3" json:"task,omitempty"`
	Tags      []string `protobuf:"bytes,3,rep,name=tags,proto3" json:"tags,omitempty"`
	Timestamp string   `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *StartParameters) Reset() {
	*x = StartParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartParameters) ProtoMessage() {}

func (x *StartParameters) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartParameters.ProtoReflect.Descriptor instead.
func (*StartParameters) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{1}
}

func (x *StartParameters) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *StartParameters) GetTask() string {
	if x != nil {
		return x.Task
	}
	return ""
}

func (x *StartParameters) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *StartParameters) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

type StopParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp string `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *StopParameters) Reset() {
	*x = StopParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopParameters) ProtoMessage() {}

func (x *StopParameters) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopParameters.ProtoReflect.Descriptor instead.
func (*StopParameters) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{2}
}

func (x *StopParameters) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

type ResumeParameters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp string `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ResumeParameters) Reset() {
	*x = ResumeParameters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResumeParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResumeParameters) ProtoMessage() {}

func (x *ResumeParameters) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResumeParameters.ProtoReflect.Descriptor instead.
func (*ResumeParameters) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{3}
}

func (x *ResumeParameters) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

type SimpleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SimpleResponse) Reset() {
	*x = SimpleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SimpleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SimpleResponse) ProtoMessage() {}

func (x *SimpleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SimpleResponse.ProtoReflect.Descriptor instead.
func (*SimpleResponse) Descriptor() ([]byte, []int) {
	return file_grpc_proto_rawDescGZIP(), []int{4}
}

func (x *SimpleResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SimpleResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_grpc_proto protoreflect.FileDescriptor

var file_grpc_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x74, 0x74,
	0x22, 0x83, 0x01, 0x0a, 0x05, 0x54, 0x69, 0x6d, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x22, 0x71, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x2e, 0x0a, 0x0e, 0x53, 0x74, 0x6f,
	0x70, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x30, 0x0a, 0x10, 0x52, 0x65, 0x73,
	0x75, 0x6d, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x40, 0x0a, 0x0e, 0x53,
	0x69, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x94, 0x01,
	0x0a, 0x02, 0x54, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x72, 0x12, 0x13, 0x2e, 0x74, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x09, 0x2e, 0x74, 0x74, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x72, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x09, 0x53, 0x74, 0x6f, 0x70, 0x54, 0x69, 0x6d, 0x65,
	0x72, 0x12, 0x12, 0x2e, 0x74, 0x74, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x09, 0x2e, 0x74, 0x74, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x72,
	0x22, 0x00, 0x12, 0x30, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x72, 0x12, 0x14, 0x2e, 0x74, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x09, 0x2e, 0x74, 0x74, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x72, 0x22, 0x00, 0x42, 0x1d, 0x5a, 0x1b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x61, 0x78, 0x6d, 0x6f, 0x65, 0x68, 0x6c, 0x2f, 0x74, 0x74, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_proto_rawDescOnce sync.Once
	file_grpc_proto_rawDescData = file_grpc_proto_rawDesc
)

func file_grpc_proto_rawDescGZIP() []byte {
	file_grpc_proto_rawDescOnce.Do(func() {
		file_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_proto_rawDescData)
	})
	return file_grpc_proto_rawDescData
}

var file_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_grpc_proto_goTypes = []interface{}{
	(*Timer)(nil),            // 0: tt.Timer
	(*StartParameters)(nil),  // 1: tt.StartParameters
	(*StopParameters)(nil),   // 2: tt.StopParameters
	(*ResumeParameters)(nil), // 3: tt.ResumeParameters
	(*SimpleResponse)(nil),   // 4: tt.SimpleResponse
}
var file_grpc_proto_depIdxs = []int32{
	1, // 0: tt.Tt.StartTimer:input_type -> tt.StartParameters
	2, // 1: tt.Tt.StopTimer:input_type -> tt.StopParameters
	3, // 2: tt.Tt.ResumeTimer:input_type -> tt.ResumeParameters
	0, // 3: tt.Tt.StartTimer:output_type -> tt.Timer
	0, // 4: tt.Tt.StopTimer:output_type -> tt.Timer
	0, // 5: tt.Tt.ResumeTimer:output_type -> tt.Timer
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_proto_init() }
func file_grpc_proto_init() {
	if File_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Timer); i {
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
		file_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartParameters); i {
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
		file_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopParameters); i {
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
		file_grpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResumeParameters); i {
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
		file_grpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SimpleResponse); i {
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
			RawDescriptor: file_grpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_proto_goTypes,
		DependencyIndexes: file_grpc_proto_depIdxs,
		MessageInfos:      file_grpc_proto_msgTypes,
	}.Build()
	File_grpc_proto = out.File
	file_grpc_proto_rawDesc = nil
	file_grpc_proto_goTypes = nil
	file_grpc_proto_depIdxs = nil
}
