
build-grpc:
	protoc --go_out=. --go-grpc_out=. --proto_path=. writy.proto

