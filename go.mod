module gitlab.uisee.ai/cloud/sdk/go2sky-grpc-plugin

go 1.14

replace gitlab.uisee.ai/cloud/sdk/utclogger => gitlab.uisee.ai/cloud/sdk/utclogger.git v0.1.1

require (
	github.com/SkyAPM/go2sky v0.5.0
	gitlab.uisee.ai/cloud/sdk/utclogger v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.33.2
	gorm.io/gorm v1.20.7
)
