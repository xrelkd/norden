.PHONY: build run clean

PROGRAM_FILE=norden
BUILD_DIR=build
VERSION=0.1.0

build:
	@go build \
		-ldflags "-s -w -X github.com/xrelkd/norden/pkg/version.Version=${VERSION}" \
	 	-o ${BUILD_DIR}/${PROGRAM_FILE}\
		cmd/norden/norden.go

run: build
	@./${BUILD_DIR}/${PROGRAM_FILE}

clean:
	@go clean
	@rm -rf ${BUILD_DIR}
