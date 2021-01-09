.PHONY: all
all:
	protoc -I/usr/local/include -I. -I${GOPATH}/src  --go_out=./ ./src/models/*.proto