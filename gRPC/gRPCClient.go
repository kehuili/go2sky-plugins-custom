package grpcPlugin

import (
	"context"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	agent "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const componentIDGOGrpcClient = 5002

func GrpcClientMiddleware(tracer *go2sky.Tracer, host string) func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) (err error) {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Logic before invoking the invoker
		var traceCtx context.Context
		span, err := tracer.CreateExitSpan(ctx, method, host, func(header string) error {
			traceCtx = metadata.AppendToOutgoingContext(ctx, propagation.Header, header)
			return nil
		})
		span.SetComponent(componentIDGOGrpcClient)
		span.Tag(go2sky.TagURL, method)
		span.SetSpanLayer(agent.SpanLayer_RPCFramework)
		defer span.End()
		// Calls the invoker to execute RPC
		err = invoker(traceCtx, method, req, reply, cc, opts...)
		return err
	}
}
