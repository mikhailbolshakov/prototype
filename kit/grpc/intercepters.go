package grpc

import (
	"context"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// this middleware is applied on server side
// it retrieves gRPC metadata and puts it to the context
func ContextUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = kitContext.FromGrpcMD(ctx, md)
		}
		resp, err := handler(ctx, req)

		return resp, err
	}
}

// this middleware is applied on server side
// it retrieves gRPC metadata and puts it to the context
func ContextStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		return err
	}
}

// this middleware is applied on client side
// it retrieves session params from the context (normally it's populated in HTTP middleware or by another caller) and puts it to gRPS metadata
func ContextUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx := context.Background()
		if md, ok := kitContext.FromContextToGrpcMD(parentCtx); ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func ContextStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(parentCtx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx := context.Background()
		if md, ok := kitContext.FromContextToGrpcMD(parentCtx); ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}