VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1.0")
COMMIT:=$(shell git rev-parse --short HEAD)
LDFLAGS:="-X github.com/filipenos/projects/cmd.Version=${VERSION} -X github.com/filipenos/projects/cmd.Commit=${COMMIT} "

clean:
	rm -rf projects

install:
	go install -ldflags=${LDFLAGS}

build:
	go build -ldflags=${LDFLAGS}
