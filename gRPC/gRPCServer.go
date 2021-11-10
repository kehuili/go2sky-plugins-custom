package grpcPlugin

import (
	"context"
	"time"

	"github.com/SkyAPM/go2sky"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentIDGOGrpcServer = 23

func GrpcServerMiddleware(tracer *go2sky.Tracer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if tracer == nil {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			r, err := handler(ctx, req)
			return r, err
		}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		span, traceCtx, err := tracer.CreateEntrySpan(ctx, info.FullMethod, func(key string) (string, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return "", nil
			}
			sw8, ok := md[key]
			if !ok {
				return "", nil
			}
			return sw8[0], nil
		})
		span.SetComponent(componentIDGOGrpcServer)
		span.Tag(go2sky.TagURL, info.FullMethod)
		span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

		//处理请求
		r, err := handler(traceCtx, req)
		if err != nil {
			span.Error(time.Now(), err.Error())
		}

		span.End()

		return r, err
	}
}
