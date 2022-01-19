export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: robot-client robot-pong robot-server

fmt:
	go fmt ./...

proto:
	protoc --go_out=pkg/net/proto pkg/net/proto/*.proto 

robot-client:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/robot-client ./cmd/robot-client

robot-pong:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/robot-pong ./cmd/robot-pong

robot-server:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/robot-server ./cmd/robot-server

clean:
	rm -f ./bin/robot-client
	rm -f ./bin/robot-pong
	rm -f ./bin/robot-server
