package proto

//go:generate protoc  --go_out=..  --go-grpc_out=.. --proto_path=. auth.proto
//go:generate protoc  --go_out=..  --go-grpc_out=.. --proto_path=. user.proto
//go:generate protoc  --go_out=..  --go-grpc_out=.. --proto_path=. common.proto
