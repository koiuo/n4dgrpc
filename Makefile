L5D_PROTOS_REPO=github.com/linkerd/linkerd/mesh/core/src/main/protobuf
PROTOS=$(GOPATH)/src/${L5D_PROTOS_REPO}
GRPC_STUBS=$(patsubst %.proto,%.pb.go,${PROTOS}/*.proto)
DOCKERTAG=m1w
BINARY=n4dgrpc

${BINARY}: ${GRPC_STUBS}
	go build -o ${BINARY}

${PROTOS}/%.pb.go: ${PROTOS}
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	./protoc -I ${PROTOS} ${PROTOS}/*.proto --go_out=plugins=grpc:${PROTOS}

${PROTOS}:
	go get -u -d ${L5D_PROTOS_REPO} || true

.PHONY: clean
clean:
	rm ${BINARY}

