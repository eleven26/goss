.PHONY: init
init:
	go build -modfile=tools/go.mod -o bin/gofumpt mvdan.cc/gofumpt
	go build -modfile=tools/go.mod -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: check
check:
	bin/golangci-lint run

FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: format
format:
	go mod tidy
	bin/gofumpt -w $(FILES)

.PHONY: test
test:
	go test -v ./core/* -cover

.PHONY: integration
integration:
	go test -v ./drivers/aliyun -cover -tags=integration
	go test -v ./drivers/tencent -cover -tags=integration
	go test -v ./drivers/qiniu -cover -tags=integration
	go test -v ./drivers/huawei -cover -tags=integration
	go test -v ./goss/* -cover -tags=integration

.PHONY: all
all:
	make check
	make format
	make test
	make integration