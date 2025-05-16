BUILD_DIRECTORY=_build
PROGRAM_NAME=bolt-ui

all: test lint build

ci: tools dependencies generate fmt check-repository-unchanged test lint build

build-directory:
	mkdir -p ./${BUILD_DIRECTORY}

build: build-directory
	go build -o ./${BUILD_DIRECTORY}/${PROGRAM_NAME} ./cmd/${PROGRAM_NAME}

build-race: build-directory
	go build -race -o ./${BUILD_DIRECTORY}/${PROGRAM_NAME} ./cmd/${PROGRAM_NAME}

frontend:
	./_tools/build_frontend.sh

check-repository-unchanged:
	./_tools/check_repository_unchanged.sh

tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/rinchsan/gosimports/cmd/gosimports@v0.3.5 # https://github.com/golang/go/issues/20818

dependencies:
	go get ./...

generate:
	 go generate ./...

lint: 
	go vet ./...
	staticcheck ./...

test:
	go test ./...

clean:
	rm -rf ./${BUILD_DIRECTORY}

fmt:
	gosimports -l -w ./

.PHONY: all build build-directory frontend check-repository-unchanged build-race tools dependencies generate lint test clean fmt
