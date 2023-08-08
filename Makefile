LDFLAGS := -s -w
BIN_NAME := "vanity"

build:
	@env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_NAME)

fmt:
	go fmt ./...

vet:
	go vet ./...

.PHONY: build fmt vet