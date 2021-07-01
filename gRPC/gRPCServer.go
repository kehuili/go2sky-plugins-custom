package grpcPlugin

import (
	"context"

	"github.com/SkyAPM/go2sky"
	"gitlab.uisee.ai/cloud/sdk/utclogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentIDGOGrpcServer = 23

func GrpcServerMiddleware(logger *utclogger.Logger, tracer *go2sky.Tracer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span, traceCtx, err := tracer.CreateEntrySpan(ctx, info.FullMethod, func(key string) (string, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				logger.Infof(nil, nil, "Retrieving metadata failed")
				return "", nil
			}
			sw8, ok := md[key]
			if !ok {
				logger.Infof(nil, nil, "sw8 is not supplied")
				return "", nil
			}
			return sw8[0], nil
		})
		span.SetComponent(componentIDGOGrpcServer)
		span.Tag(go2sky.TagURL, info.FullMethod)
		span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

		//处理请求
		r, err := handler(traceCtx, req)

		span.End()

		return r, err
	}
}
