protoc --proto_path=api/proto/v1 --go-grpc_out=pkg/api/v1 todo-service.proto
protoc --proto_path=api/proto/v1 --go_out=require_unimplemented_servers=false:pkg/api/v1 todo-service.proto
