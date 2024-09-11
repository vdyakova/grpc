LOCAL_BIN:=$(CURDIR)/bin

install-deps:

	set GOBIN=D:/work/microservice_course/hw1/bin && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	set	GOBIN=D:/work/microservice_course/hw1/bin && go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install  google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


generate:
	make generate-note-api

generate-note-api:

	protoc --proto_path=api/note_v1 \
		   --go_out=pkg/note_v1 --go_opt=paths=source_relative \
 		  --go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
 		  api/note_v1/note.proto
build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/auth.go

copy-to-server:
	scp service_linux root@31.41.154.33:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/promise/test-server:v0.0.1 .
	docker login -u token -p CRgAAAAATm9NHHvPmdiRcFiX22NeS-h9ieBhPDH0  cr.selcloud.ru/promise
	docker push cr.selcloud.ru/promise/test-server:v0.0.1
