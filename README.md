# Go2sky for gRPC

## Installation

```bash
go get -u gitlab.uisee.ai/cloud/sdk/go2sky-grpc-plugin
```

## Usage
```go
package main

import (
	"log"

	"github.com/SkyAPM/go2sky"
	grpcPlugin "gitlab.uisee.ai/cloud/sdk/go2sky-grpc-plugin"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
)

func main() {
  rp, err := reporter.NewGRPCReporter("skywalking-uri")
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer re.Close()

	tracer, err := go2sky.NewTracer("grpc-server", go2sky.WithReporter(re))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}

	logger := utclogger.New("console-xxx")

	// Use grpc server middleware with tracing
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcPlugin.GrpcServerMiddleware(logger, tracer)))

  // do something
  
	// Use grpc client middleware with tracing
  conn, _ := grpc.Dial("dest", grpc.WithInsecure(), grpc.WithUnaryInterceptor(GrpcClientMiddleware(tracer, "localhost")))
  
  // do something
}
```
