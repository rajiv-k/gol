GOOS ?= linux
BINARY_NAME ?= gol

gol: main.go
	GOOS=${GOOS} go build -o ${BINARY_NAME} main.go

test:
	go test -v ./...

clean:
	@rm -rf ${BINARY_NAME}

.PHONY: gol clean test
