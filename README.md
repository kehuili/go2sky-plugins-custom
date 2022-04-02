# Go2sky for gRPC

## Installation

```bash
go get -u gitlab.uisee.ai/cloud/sdk/go2sky-plugin
```

## Usage
```go
package main

import (
	"log"

	"github.com/SkyAPM/go2sky"
	mqttPlugin "gitlab.uisee.ai/cloud/sdk/go2sky-plugin/paho"
	grpcPlugin "gitlab.uisee.ai/cloud/sdk/go2sky-plugin/gRPC"
	gormPlugin "gitlab.uisee.ai/cloud/sdk/go2sky-plugin/gorm"
	ginPlugin "gitlab.uisee.ai/cloud/sdk/go2sky-plugin/gin"
	skyUtil "gitlab.uisee.ai/cloud/sdk/go2sky-plugin/util"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
  rp, tracer, err := skyUtil.Reporter("skywalking-uri", "server-name")
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	defer skyUtil.Close(rp)

	// Use grpc server middleware with tracing
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcPlugin.GrpcServerMiddleware(tracer)))
  
	// Use grpc client middleware with tracing
  conn, _ := grpc.Dial("dest", grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpcPlugin.GrpcClientMiddleware(tracer, "localhost")))

  // gorm plugin
  db, _ = gorm.Open(mysql.Open(dbDsn), &gorm.Config{})
  gormPlugin.RegisterAll(db, tracer, dbDsn, gormPlugin.GormCallback)
  db.WithContext(ctx).Create()

  // mqtt publish
	text := fmt.Sprintf(`{"a": "Message hello"}`)
	msg, err := mqttPlugin.BeforePublish(tracer, "test", topic, text, ctx)
	client.Publish(topic, 0, false, msg)

  // mqtt receive
  ctx := context.Background()
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		ctx,err = mqttPlugin.AfterOnMessage(tracer, topic, msg.Payload(), ctx)
	}

  //gin
	router := gin.New()
	router.Use(ginPlugin.Middleware(router, tracer))
  // exclude ping
	router.Use(ginPlugin.Middleware(router, tracer, ginPlugin.WithExcludePaths([]string{"/ping"})))
  // for some url, read sw8 from body
	router.Use(ginPlugin.Middleware(router, tracer, ginPlugin.WithFromBodyPaths([]string{"/notification"})))
}
```
