L5D_PROTOS_REPO=github.com/linkerd/linkerd/mesh/core/src/main/protobuf
PROTOS=$(GOPATH)/src/${L5D_PROTOS_REPO}
GRPC_STUBS=$(patsubst %.proto,%.pb.go,${PROTOS}/*.proto)
DOCKERTAG=m1w
BINARY=n4dgrpc

default: all
.PHONY: all
all: deps ${BINARY}

${BINARY}:
	go build -o ${BINARY}

${PROTOS}/%.pb.go: ${PROTOS}
	go get -u github.com/golang/protobuf/protoc-gen-go
	./protoc -I ${PROTOS} ${PROTOS}/*.proto --go_out=plugins=grpc:${PROTOS}

${PROTOS}:
	go get -u -d ${L5D_PROTOS_REPO} || true

.PHONY: test
test:
	go test -v ./...

.PHONY: deps
deps: ${GRPC_STUBS}
	go get -u ./...

.PHONY: test-deps
test-deps: ${GRPC_STUBS}
	go get -t -u ./...

.PHONY: clean
clean:
	rm ${BINARY}

