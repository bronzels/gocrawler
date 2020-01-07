PROTO_FILES := $(shell find ../pb -name "*.proto")
PROTO_TARGETS := $(PROTO_FILES:.proto=.pb.go)

pb: $(PROTO_TARGETS)

%.pb.go: %.proto
	protoc -I ../pb $< --go_out=plugins=grpc:./pb

build:
	#go build ./src/helloworldserver
	#go build ./src/helloworldclient
	go build ./src/crawlnewsserver
	go build ./src/crawlnewsclient

test:
	go test -cover ./test/*.go

all: pb build