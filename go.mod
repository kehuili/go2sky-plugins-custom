module gitlab.uisee.ai/cloud/sdk/go2sky-plugin

go 1.14

replace gitlab.uisee.ai/cloud/sdk/utclogger => gitlab.uisee.ai/cloud/sdk/utclogger.git v0.1.1

require (
	github.com/SkyAPM/go2sky v1.1.0
	github.com/gin-gonic/gin v1.6.3
	gitlab.uisee.ai/cloud/sdk/utclogger v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.38.0
	gorm.io/gorm v1.20.7
	skywalking.apache.org/repo/goapi v0.0.0-20210628073857-a95ba03d3c7a
)
