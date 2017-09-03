L5D_PROTOS_REPO=github.com/linkerd/linkerd/mesh/core/src/main/protobuf
PROTOS=$(GOPATH)/src/${L5D_PROTOS_REPO}
GRPC_STUBS=$(patsubst %.proto,%.pb.go,${PROTOS}/*.proto)
DOCKERTAG=m1w
BINARY=n4dgrpc
VERSION=$(shell git describe 2>/dev/null || echo "0.0.0")
RELEASE=release
RELEASE_FORMAT=${BINARY}-${GOOS}-${GOARCH}-${VERSION}.tar.gz

default: bin
.PHONY: bin
bin: deps ${BINARY}

.PHONY: release
release: release/n4dgrpc-linux-amd64-${VERSION}.tar.gz
release: release/n4dgrpc-linux-386-${VERSION}.tar.gz
release: release/n4dgrpc-darwin-amd64-${VERSION}.tar.gz

.PHONY: release-dir
release-dir:
	mkdir -p ${RELEASE}

${RELEASE}/%-${VERSION}.tar.gz: | release-dir
	GOOS=$(shell echo $* | cut -d '-' -f 2) GOARCH=$(shell echo $* | cut -d '-' -f 3) \
	go build -ldflags="-s -w" -o ${RELEASE}/${BINARY}
	tar -C ${RELEASE} -czf $@ ${BINARY}
	rm ${RELEASE}/${BINARY}

${BINARY}:
	go build -o $@

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
	rm -f ${BINARY}
	rm -rf ${RELEASE}
