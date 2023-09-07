// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: coordinate.proto

package coordinate

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

type WorkerStatus int32

const (
	WorkerStatus_RunningStatus WorkerStatus = 0
	WorkerStatus_OfflineStatus WorkerStatus = 1
)

// Enum value maps for WorkerStatus.
var (
	WorkerStatus_name = map[int32]string{
		0: "RunningStatus",
		1: "OfflineStatus",
	}
	WorkerStatus_value = map[string]int32{
		"RunningStatus": 0,
		"OfflineStatus": 1,
	}
)

func (x WorkerStatus) Enum() *WorkerStatus {
	p := new(WorkerStatus)
	*p = x
	return p
}

func (x WorkerStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WorkerStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_coordinate_proto_enumTypes[0].Descriptor()
}

func (WorkerStatus) Type() protoreflect.EnumType {
	return &file_coordinate_proto_enumTypes[0]
}

func (x WorkerStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WorkerStatus.Descriptor instead.
func (WorkerStatus) EnumDescriptor() ([]byte, []int) {
	return file_coordinate_proto_rawDescGZIP(), []int{0}
}

type TaskStatus int32

const (
	TaskStatus_TaskRunningStatus  TaskStatus = 0
	TaskStatus_TaskFinishedStatus TaskStatus = 1
	TaskStatus_TaskFailedStatus   TaskStatus = 2
)

// Enum value maps for TaskStatus.
var (
	TaskStatus_name = map[int32]string{
		0: "TaskRunningStatus",
		1: "TaskFinishedStatus",
		2: "TaskFailedStatus",
	}
	TaskStatus_value = map[string]int32{
		"TaskRunningStatus":  0,
		"TaskFinishedStatus": 1,
		"TaskFailedStatus":   2,
	}
)

func (x TaskStatus) Enum() *TaskStatus {
	p := new(TaskStatus)
	*p = x
	return p
}

func (x TaskStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_coordinate_proto_enumTypes[1].Descriptor()
}

func (TaskStatus) Type() protoreflect.EnumType {
	return &file_coordinate_proto_enumTypes[1]
}

func (x TaskStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskStatus.Descriptor instead.
func (TaskStatus) EnumDescriptor() ([]byte, []int) {
	return file_coordinate_proto_rawDescGZIP(), []int{1}
}

type ClientMode int32

const (
	ClientMode_MapMode    ClientMode = 0
	ClientMode_ReduceMode ClientMode = 1
)

// Enum value maps for ClientMode.
var (
	ClientMode_name = map[int32]string{
		0: "MapMode",
		1: "ReduceMode",
	}
	ClientMode_value = map[string]int32{
		"MapMode":    0,
		"ReduceMode": 1,
	}
)

func (x ClientMode) Enum() *ClientMode {
	p := new(ClientMode)
	*p = x
	return p
}

func (x ClientMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientMode) Descriptor() protoreflect.EnumDescriptor {
	return file_coordinate_proto_enumTypes[2].Descriptor()
}

func (ClientMode) Type() protoreflect.EnumType {
	return &file_coordinate_proto_enumTypes[2]
}

func (x ClientMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClientMode.Descriptor instead.
func (ClientMode) EnumDescriptor() ([]byte, []int) {
	return file_coordinate_proto_rawDescGZIP(), []int{2}
}

type ClientInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address string     `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Mode    ClientMode `protobuf:"varint,3,opt,name=mode,proto3,enum=coordinate.ClientMode" json:"mode,omitempty"`
}

func (x *ClientInfo) Reset() {
	*x = ClientInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientInfo) ProtoMessage() {}

func (x *ClientInfo) ProtoReflect() protoreflect.Message {
	mi := &file_coordinate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientInfo.ProtoReflect.Descriptor instead.
func (*ClientInfo) Descriptor() ([]byte, []int) {
	return file_coordinate_proto_rawDescGZIP(), []int{0}
}

func (x *ClientInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ClientInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *ClientInfo) GetMode() ClientMode {
	if x != nil {
		return x.Mode
	}
	return ClientMode_MapMode
}

type HeartbeatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 服务当前状态
	Name        string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status      WorkerStatus `protobuf:"varint,2,opt,name=status,proto3,enum=coordinate.WorkerStatus" json:"status,omitempty"`
	CurrentTask string       `protobuf:"bytes,3,opt,name=current_task,json=currentTask,proto3" json:"current_task,omitempty"`
	TaskStatus  TaskStatus   `protobuf:"varint,4,opt,name=task_status,json=taskStatus,proto3,enum=coordinate.TaskStatus" json:"task_status,omitempty"`
	Filenames   []string     `protobuf:"bytes,5,rep,name=filenames,proto3" json:"filenames,omitempty"`
	Message     string       `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *HeartbeatRequest) Reset() {
	*x = HeartbeatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatRequest) ProtoMessage() {}

func (x *HeartbeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_coordinate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatRequest.ProtoReflect.Descriptor instead.
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return file_coordinate_proto_rawDescGZIP(), []int{1}
}

func (x *HeartbeatRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HeartbeatRequest) GetStatus() WorkerStatus {
	if x != nil {
		return x.Status
	}
	return WorkerStatus_RunningStatus
}

func (x *HeartbeatRequest) GetCurrentTask() string {
	if x != nil {
		return x.CurrentTask
	}
	return ""
}

func (x *HeartbeatRequest) GetTaskStatus() TaskStatus {
	if x != nil {
		return x.TaskStatus
	}
	return TaskStatus_TaskRunningStatus
}

func (x *HeartbeatRequest) GetFilenames() []string {
	if x != nil {
		return x.Filenames
	}
	return nil
}

func (x *HeartbeatRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coordinate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_coordinate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_coordinate_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_coordinate_proto protoreflect.FileDescriptor

var file_coordinate_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x22, 0x66,
	0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2a, 0x0a, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64,
	0x69, 0x6e, 0x61, 0x74, 0x65, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65,
	0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x22, 0xec, 0x01, 0x0a, 0x10, 0x48, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x30, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x18, 0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x2e, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x61, 0x73,
	0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x54, 0x61, 0x73, 0x6b, 0x12, 0x37, 0x0a, 0x0b, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6f, 0x72,
	0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x38, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a,
	0x34, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x11, 0x0a, 0x0d, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x10, 0x01, 0x2a, 0x51, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x75, 0x6e, 0x6e, 0x69,
	0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x54, 0x61,
	0x73, 0x6b, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x61, 0x73, 0x6b, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x10, 0x02, 0x2a, 0x29, 0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x61, 0x70, 0x4d, 0x6f, 0x64,
	0x65, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x4d, 0x6f, 0x64,
	0x65, 0x10, 0x01, 0x32, 0x87, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61,
	0x74, 0x65, 0x12, 0x38, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x16,
	0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x2e, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e,
	0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x09,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x1c, 0x2e, 0x63, 0x6f, 0x6f, 0x72,
	0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69,
	0x6e, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a,
	0x0b, 0x2f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coordinate_proto_rawDescOnce sync.Once
	file_coordinate_proto_rawDescData = file_coordinate_proto_rawDesc
)

func file_coordinate_proto_rawDescGZIP() []byte {
	file_coordinate_proto_rawDescOnce.Do(func() {
		file_coordinate_proto_rawDescData = protoimpl.X.CompressGZIP(file_coordinate_proto_rawDescData)
	})
	return file_coordinate_proto_rawDescData
}

var file_coordinate_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_coordinate_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_coordinate_proto_goTypes = []interface{}{
	(WorkerStatus)(0),        // 0: coordinate.WorkerStatus
	(TaskStatus)(0),          // 1: coordinate.TaskStatus
	(ClientMode)(0),          // 2: coordinate.ClientMode
	(*ClientInfo)(nil),       // 3: coordinate.ClientInfo
	(*HeartbeatRequest)(nil), // 4: coordinate.HeartbeatRequest
	(*Response)(nil),         // 5: coordinate.Response
}
var file_coordinate_proto_depIdxs = []int32{
	2, // 0: coordinate.ClientInfo.mode:type_name -> coordinate.ClientMode
	0, // 1: coordinate.HeartbeatRequest.status:type_name -> coordinate.WorkerStatus
	1, // 2: coordinate.HeartbeatRequest.task_status:type_name -> coordinate.TaskStatus
	3, // 3: coordinate.Coordinate.Register:input_type -> coordinate.ClientInfo
	4, // 4: coordinate.Coordinate.Heartbeat:input_type -> coordinate.HeartbeatRequest
	5, // 5: coordinate.Coordinate.Register:output_type -> coordinate.Response
	5, // 6: coordinate.Coordinate.Heartbeat:output_type -> coordinate.Response
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_coordinate_proto_init() }
func file_coordinate_proto_init() {
	if File_coordinate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coordinate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientInfo); i {
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
		file_coordinate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatRequest); i {
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
		file_coordinate_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_coordinate_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_coordinate_proto_goTypes,
		DependencyIndexes: file_coordinate_proto_depIdxs,
		EnumInfos:         file_coordinate_proto_enumTypes,
		MessageInfos:      file_coordinate_proto_msgTypes,
	}.Build()
	File_coordinate_proto = out.File
	file_coordinate_proto_rawDesc = nil
	file_coordinate_proto_goTypes = nil
	file_coordinate_proto_depIdxs = nil
}
