//*
// Aembit Edge to Aembit Cloud communication-related messages/services
//
//

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: aembit.proto

package provider

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EdgeCommander_GetCommands_FullMethodName             = "/aembit.EdgeCommander/GetCommands"
	EdgeCommander_GetConfiguration_FullMethodName        = "/aembit.EdgeCommander/GetConfiguration"
	EdgeCommander_GetPolicy_FullMethodName               = "/aembit.EdgeCommander/GetPolicy"
	EdgeCommander_GetCredential_FullMethodName           = "/aembit.EdgeCommander/GetCredential"
	EdgeCommander_GetCredentials_FullMethodName          = "/aembit.EdgeCommander/GetCredentials"
	EdgeCommander_GetCertificate_FullMethodName          = "/aembit.EdgeCommander/GetCertificate"
	EdgeCommander_ReportEvent_FullMethodName             = "/aembit.EdgeCommander/ReportEvent"
	EdgeCommander_ReportEvents_FullMethodName            = "/aembit.EdgeCommander/ReportEvents"
	EdgeCommander_RegisterAgentController_FullMethodName = "/aembit.EdgeCommander/RegisterAgentController"
)

// EdgeCommanderClient is the client API for EdgeCommander service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EdgeCommanderClient interface {
	// The long poll API called by a client to wait until the backend has command
	// It should be called once per agent (even if an agent serves multiple wokloads)
	GetCommands(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (EdgeCommander_GetCommandsClient, error)
	// Get the base configuration
	GetConfiguration(ctx context.Context, in *ConfigurationRequest, opts ...grpc.CallOption) (*ConfigurationResponse, error)
	// Get a dynamic policy
	GetPolicy(ctx context.Context, in *PolicyRequest, opts ...grpc.CallOption) (*PolicyResponse, error)
	// Get credential for identified target workload
	GetCredential(ctx context.Context, in *CredentialRequest, opts ...grpc.CallOption) (*CredentialResponse, error)
	// Get credentials for identified target workload
	GetCredentials(ctx context.Context, in *CredentialsRequest, opts ...grpc.CallOption) (*CredentialsResponse, error)
	// Get a certificate to terminate a client's TLS connection.
	GetCertificate(ctx context.Context, in *CertificateRequest, opts ...grpc.CallOption) (*CertificateResponse, error)
	ReportEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*EventResponse, error)
	ReportEvents(ctx context.Context, in *EventRequests, opts ...grpc.CallOption) (*EventResponse, error)
	RegisterAgentController(ctx context.Context, in *AgentControllerRegistrationRequest, opts ...grpc.CallOption) (*AgentControllerRegistrationResponse, error)
}

type edgeCommanderClient struct {
	cc grpc.ClientConnInterface
}

func NewEdgeCommanderClient(cc grpc.ClientConnInterface) EdgeCommanderClient {
	return &edgeCommanderClient{cc}
}

func (c *edgeCommanderClient) GetCommands(ctx context.Context, in *CommandRequest, opts ...grpc.CallOption) (EdgeCommander_GetCommandsClient, error) {
	stream, err := c.cc.NewStream(ctx, &EdgeCommander_ServiceDesc.Streams[0], EdgeCommander_GetCommands_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &edgeCommanderGetCommandsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EdgeCommander_GetCommandsClient interface {
	Recv() (*CommandResponse, error)
	grpc.ClientStream
}

type edgeCommanderGetCommandsClient struct {
	grpc.ClientStream
}

func (x *edgeCommanderGetCommandsClient) Recv() (*CommandResponse, error) {
	m := new(CommandResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *edgeCommanderClient) GetConfiguration(ctx context.Context, in *ConfigurationRequest, opts ...grpc.CallOption) (*ConfigurationResponse, error) {
	out := new(ConfigurationResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_GetConfiguration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) GetPolicy(ctx context.Context, in *PolicyRequest, opts ...grpc.CallOption) (*PolicyResponse, error) {
	out := new(PolicyResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_GetPolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) GetCredential(ctx context.Context, in *CredentialRequest, opts ...grpc.CallOption) (*CredentialResponse, error) {
	out := new(CredentialResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_GetCredential_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) GetCredentials(ctx context.Context, in *CredentialsRequest, opts ...grpc.CallOption) (*CredentialsResponse, error) {
	out := new(CredentialsResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_GetCredentials_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) GetCertificate(ctx context.Context, in *CertificateRequest, opts ...grpc.CallOption) (*CertificateResponse, error) {
	out := new(CertificateResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_GetCertificate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) ReportEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_ReportEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) ReportEvents(ctx context.Context, in *EventRequests, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_ReportEvents_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeCommanderClient) RegisterAgentController(ctx context.Context, in *AgentControllerRegistrationRequest, opts ...grpc.CallOption) (*AgentControllerRegistrationResponse, error) {
	out := new(AgentControllerRegistrationResponse)
	err := c.cc.Invoke(ctx, EdgeCommander_RegisterAgentController_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EdgeCommanderServer is the server API for EdgeCommander service.
// All implementations must embed UnimplementedEdgeCommanderServer
// for forward compatibility
type EdgeCommanderServer interface {
	// The long poll API called by a client to wait until the backend has command
	// It should be called once per agent (even if an agent serves multiple wokloads)
	GetCommands(*CommandRequest, EdgeCommander_GetCommandsServer) error
	// Get the base configuration
	GetConfiguration(context.Context, *ConfigurationRequest) (*ConfigurationResponse, error)
	// Get a dynamic policy
	GetPolicy(context.Context, *PolicyRequest) (*PolicyResponse, error)
	// Get credential for identified target workload
	GetCredential(context.Context, *CredentialRequest) (*CredentialResponse, error)
	// Get credentials for identified target workload
	GetCredentials(context.Context, *CredentialsRequest) (*CredentialsResponse, error)
	// Get a certificate to terminate a client's TLS connection.
	GetCertificate(context.Context, *CertificateRequest) (*CertificateResponse, error)
	ReportEvent(context.Context, *EventRequest) (*EventResponse, error)
	ReportEvents(context.Context, *EventRequests) (*EventResponse, error)
	RegisterAgentController(context.Context, *AgentControllerRegistrationRequest) (*AgentControllerRegistrationResponse, error)
	mustEmbedUnimplementedEdgeCommanderServer()
}

// UnimplementedEdgeCommanderServer must be embedded to have forward compatible implementations.
type UnimplementedEdgeCommanderServer struct {
}

func (UnimplementedEdgeCommanderServer) GetCommands(*CommandRequest, EdgeCommander_GetCommandsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetCommands not implemented")
}
func (UnimplementedEdgeCommanderServer) GetConfiguration(context.Context, *ConfigurationRequest) (*ConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfiguration not implemented")
}
func (UnimplementedEdgeCommanderServer) GetPolicy(context.Context, *PolicyRequest) (*PolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPolicy not implemented")
}
func (UnimplementedEdgeCommanderServer) GetCredential(context.Context, *CredentialRequest) (*CredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCredential not implemented")
}
func (UnimplementedEdgeCommanderServer) GetCredentials(context.Context, *CredentialsRequest) (*CredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCredentials not implemented")
}
func (UnimplementedEdgeCommanderServer) GetCertificate(context.Context, *CertificateRequest) (*CertificateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCertificate not implemented")
}
func (UnimplementedEdgeCommanderServer) ReportEvent(context.Context, *EventRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportEvent not implemented")
}
func (UnimplementedEdgeCommanderServer) ReportEvents(context.Context, *EventRequests) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportEvents not implemented")
}
func (UnimplementedEdgeCommanderServer) RegisterAgentController(context.Context, *AgentControllerRegistrationRequest) (*AgentControllerRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAgentController not implemented")
}
func (UnimplementedEdgeCommanderServer) mustEmbedUnimplementedEdgeCommanderServer() {}

// UnsafeEdgeCommanderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EdgeCommanderServer will
// result in compilation errors.
type UnsafeEdgeCommanderServer interface {
	mustEmbedUnimplementedEdgeCommanderServer()
}

func RegisterEdgeCommanderServer(s grpc.ServiceRegistrar, srv EdgeCommanderServer) {
	s.RegisterService(&EdgeCommander_ServiceDesc, srv)
}

func _EdgeCommander_GetCommands_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CommandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EdgeCommanderServer).GetCommands(m, &edgeCommanderGetCommandsServer{stream})
}

type EdgeCommander_GetCommandsServer interface {
	Send(*CommandResponse) error
	grpc.ServerStream
}

type edgeCommanderGetCommandsServer struct {
	grpc.ServerStream
}

func (x *edgeCommanderGetCommandsServer) Send(m *CommandResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _EdgeCommander_GetConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).GetConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_GetConfiguration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).GetConfiguration(ctx, req.(*ConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_GetPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).GetPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_GetPolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).GetPolicy(ctx, req.(*PolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_GetCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).GetCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_GetCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).GetCredential(ctx, req.(*CredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_GetCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).GetCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_GetCredentials_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).GetCredentials(ctx, req.(*CredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_GetCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).GetCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_GetCertificate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).GetCertificate(ctx, req.(*CertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_ReportEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).ReportEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_ReportEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).ReportEvent(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_ReportEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequests)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).ReportEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_ReportEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).ReportEvents(ctx, req.(*EventRequests))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeCommander_RegisterAgentController_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AgentControllerRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeCommanderServer).RegisterAgentController(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeCommander_RegisterAgentController_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeCommanderServer).RegisterAgentController(ctx, req.(*AgentControllerRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EdgeCommander_ServiceDesc is the grpc.ServiceDesc for EdgeCommander service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EdgeCommander_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "aembit.EdgeCommander",
	HandlerType: (*EdgeCommanderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConfiguration",
			Handler:    _EdgeCommander_GetConfiguration_Handler,
		},
		{
			MethodName: "GetPolicy",
			Handler:    _EdgeCommander_GetPolicy_Handler,
		},
		{
			MethodName: "GetCredential",
			Handler:    _EdgeCommander_GetCredential_Handler,
		},
		{
			MethodName: "GetCredentials",
			Handler:    _EdgeCommander_GetCredentials_Handler,
		},
		{
			MethodName: "GetCertificate",
			Handler:    _EdgeCommander_GetCertificate_Handler,
		},
		{
			MethodName: "ReportEvent",
			Handler:    _EdgeCommander_ReportEvent_Handler,
		},
		{
			MethodName: "ReportEvents",
			Handler:    _EdgeCommander_ReportEvents_Handler,
		},
		{
			MethodName: "RegisterAgentController",
			Handler:    _EdgeCommander_RegisterAgentController_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetCommands",
			Handler:       _EdgeCommander_GetCommands_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "aembit.proto",
}
