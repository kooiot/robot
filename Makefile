export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: robot-client robot-server

fmt:
	go fmt ./...

proto:
	protoc --go_out=pkg/net/proto pkg/net/proto/*.proto 

robot-client:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/robot-client ./cmd/client

robot-server:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/robot-server ./cmd/server

clean:
	rm -f ./bin/robot-client
	rm -f ./bin/robot-server
