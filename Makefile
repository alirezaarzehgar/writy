IMG := alirezaarzehgar/writy:$(shell cat VERSION)

build:
	docker build . --file Dockerfile --tag ${IMG}

push:
	docker push ${IMG}


build-grpc:
	protoc --go_out=. --go-grpc_out=. --proto_path=. writy.proto

