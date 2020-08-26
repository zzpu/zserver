module zserver

go 1.13

require (
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/google/wire v0.4.0
	github.com/gorilla/websocket v1.4.2
	github.com/hpcloud/tail v1.0.0
	github.com/zzpu/kratos v0.4.5
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	google.golang.org/grpc v1.27.0

)

replace google.golang.org/grpc => google.golang.org/grpc v1.24.0
