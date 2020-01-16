module clientmsg

go 1.12

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/pborman/uuid v1.2.0
	github.com/spf13/viper v1.5.0
	google.golang.org/grpc v1.23.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.25.1
