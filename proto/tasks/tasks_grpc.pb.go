// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package tasks

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// TasksClient is the client API for Tasks service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TasksClient interface {
	New(ctx context.Context, in *NewTaskRequest, opts ...grpc.CallOption) (*Task, error)
	NextTransitions(ctx context.Context, in *NextTransitionsRequest, opts ...grpc.CallOption) (*NextTransitionsResponse, error)
	MakeTransition(ctx context.Context, in *MakeTransitionRequest, opts ...grpc.CallOption) (*Task, error)
	GetByChannel(ctx context.Context, in *GetByChannelRequest, opts ...grpc.CallOption) (*GetByChannelResponse, error)
	SetAssignee(ctx context.Context, in *SetAssigneeRequest, opts ...grpc.CallOption) (*Task, error)
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*Task, error)
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	GetAssignmentLog(ctx context.Context, in *AssignmentLogRequest, opts ...grpc.CallOption) (*AssignmentLogResponse, error)
}

type tasksClient struct {
	cc grpc.ClientConnInterface
}

func NewTasksClient(cc grpc.ClientConnInterface) TasksClient {
	return &tasksClient{cc}
}

func (c *tasksClient) New(ctx context.Context, in *NewTaskRequest, opts ...grpc.CallOption) (*Task, error) {
	out := new(Task)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/New", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) NextTransitions(ctx context.Context, in *NextTransitionsRequest, opts ...grpc.CallOption) (*NextTransitionsResponse, error) {
	out := new(NextTransitionsResponse)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/NextTransitions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) MakeTransition(ctx context.Context, in *MakeTransitionRequest, opts ...grpc.CallOption) (*Task, error) {
	out := new(Task)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/MakeTransition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) GetByChannel(ctx context.Context, in *GetByChannelRequest, opts ...grpc.CallOption) (*GetByChannelResponse, error) {
	out := new(GetByChannelResponse)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/GetByChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) SetAssignee(ctx context.Context, in *SetAssigneeRequest, opts ...grpc.CallOption) (*Task, error) {
	out := new(Task)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/SetAssignee", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*Task, error) {
	out := new(Task)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/GetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksClient) GetAssignmentLog(ctx context.Context, in *AssignmentLogRequest, opts ...grpc.CallOption) (*AssignmentLogResponse, error) {
	out := new(AssignmentLogResponse)
	err := c.cc.Invoke(ctx, "/tasks.Tasks/GetAssignmentLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TasksServer is the server API for Tasks service.
// All implementations must embed UnimplementedTasksServer
// for forward compatibility
type TasksServer interface {
	New(context.Context, *NewTaskRequest) (*Task, error)
	NextTransitions(context.Context, *NextTransitionsRequest) (*NextTransitionsResponse, error)
	MakeTransition(context.Context, *MakeTransitionRequest) (*Task, error)
	GetByChannel(context.Context, *GetByChannelRequest) (*GetByChannelResponse, error)
	SetAssignee(context.Context, *SetAssigneeRequest) (*Task, error)
	GetById(context.Context, *GetByIdRequest) (*Task, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	GetAssignmentLog(context.Context, *AssignmentLogRequest) (*AssignmentLogResponse, error)
	mustEmbedUnimplementedTasksServer()
}

// UnimplementedTasksServer must be embedded to have forward compatible implementations.
type UnimplementedTasksServer struct {
}

func (UnimplementedTasksServer) New(context.Context, *NewTaskRequest) (*Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method New not implemented")
}
func (UnimplementedTasksServer) NextTransitions(context.Context, *NextTransitionsRequest) (*NextTransitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextTransitions not implemented")
}
func (UnimplementedTasksServer) MakeTransition(context.Context, *MakeTransitionRequest) (*Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeTransition not implemented")
}
func (UnimplementedTasksServer) GetByChannel(context.Context, *GetByChannelRequest) (*GetByChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByChannel not implemented")
}
func (UnimplementedTasksServer) SetAssignee(context.Context, *SetAssigneeRequest) (*Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAssignee not implemented")
}
func (UnimplementedTasksServer) GetById(context.Context, *GetByIdRequest) (*Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedTasksServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedTasksServer) GetAssignmentLog(context.Context, *AssignmentLogRequest) (*AssignmentLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAssignmentLog not implemented")
}
func (UnimplementedTasksServer) mustEmbedUnimplementedTasksServer() {}

// UnsafeTasksServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TasksServer will
// result in compilation errors.
type UnsafeTasksServer interface {
	mustEmbedUnimplementedTasksServer()
}

func RegisterTasksServer(s grpc.ServiceRegistrar, srv TasksServer) {
	s.RegisterService(&_Tasks_serviceDesc, srv)
}

func _Tasks_New_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).New(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/New",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).New(ctx, req.(*NewTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_NextTransitions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NextTransitionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).NextTransitions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/NextTransitions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).NextTransitions(ctx, req.(*NextTransitionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_MakeTransition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MakeTransitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).MakeTransition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/MakeTransition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).MakeTransition(ctx, req.(*MakeTransitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_GetByChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).GetByChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/GetByChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).GetByChannel(ctx, req.(*GetByChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_SetAssignee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAssigneeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).SetAssignee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/SetAssignee",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).SetAssignee(ctx, req.(*SetAssigneeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/GetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tasks_GetAssignmentLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignmentLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksServer).GetAssignmentLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.Tasks/GetAssignmentLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksServer).GetAssignmentLog(ctx, req.(*AssignmentLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tasks_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tasks.Tasks",
	HandlerType: (*TasksServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "New",
			Handler:    _Tasks_New_Handler,
		},
		{
			MethodName: "NextTransitions",
			Handler:    _Tasks_NextTransitions_Handler,
		},
		{
			MethodName: "MakeTransition",
			Handler:    _Tasks_MakeTransition_Handler,
		},
		{
			MethodName: "GetByChannel",
			Handler:    _Tasks_GetByChannel_Handler,
		},
		{
			MethodName: "SetAssignee",
			Handler:    _Tasks_SetAssignee_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _Tasks_GetById_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Tasks_Search_Handler,
		},
		{
			MethodName: "GetAssignmentLog",
			Handler:    _Tasks_GetAssignmentLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/tasks/tasks.proto",
}
