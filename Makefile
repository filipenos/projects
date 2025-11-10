VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1.0")
COMMIT:=$(shell git rev-parse --short HEAD)
LDFLAGS:="-X github.com/filipenos/projects/pkg/command.Version=${VERSION} -X github.com/filipenos/projects/pkg/command.Commit=${COMMIT} "

clean:
	rm -rf projects .cache

install:
	go install -ldflags=${LDFLAGS}

build:
	go build -ldflags=${LDFLAGS}

test:
	GOCACHE=$$(pwd)/.cache go test ./...
