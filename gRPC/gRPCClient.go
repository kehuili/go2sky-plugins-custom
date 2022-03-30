package grpcPlugin

import (
	"context"
	"time"

	"github.com/SkyAPM/go2sky"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const componentIDGOGrpcClient = 23

func GrpcClientMiddleware(tracer *go2sky.Tracer, dstPeerName string) grpc.UnaryClientInterceptor {
	// 将 skywalking header 信息放入 context
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if tracer == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		outgoingMD := metadata.MD{}
		peer := dstPeerName

		// 默认用连接的目标地址
		if peer == "" {
			peer = cc.Target()
		}

		span, newCtx, err := tracer.CreateExitSpanWithContext(ctx, method, peer, func(key, value string) error {
			outgoingMD.Set(key, value)
			return nil
		})

		// 当创建 span 失败，直接调用接口，忽略 skywalking
		if err != nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		defer span.End()

		span.SetComponent(componentIDGOGrpcClient)
		span.Tag(go2sky.TagURL, method)
		span.SetSpanLayer(agentv3.SpanLayer_RPCFramework)

		// context 塞入 outgoing metadata
		newCtx = metadata.NewOutgoingContext(newCtx, outgoingMD)

		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.Error(time.Now(), err.Error())
			return err
		}

		return nil
	}
}
