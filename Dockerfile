FROM golang:1.23.1-alpine AS builder

COPY . /github.com/vdyakova/grpc/source/
WORKDIR /github.com/vdyakova/grpc/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/auth.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/vdyakova/grpc/source/bin/crud_server .

CMD ["./crud_server"]