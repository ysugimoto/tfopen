.PHONY: test benchmark

BUILD_VERSION=$(or ${VERSION}, dev)

test:
	go test ./...

linux_amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/tfopen-linux-amd64 .
	cd ./dist/ && cp ./tfopen-linux-amd64 ./tfopen && tar cfz tfopen-linux-amd64.tar.gz ./tfopen

linux_arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/tfopen-linux-arm64 .
	cd ./dist/ && cp ./tfopen-linux-arm64 ./tfopen && tar cfz tfopen-linux-arm64.tar.gz ./tfopen

darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/tfopen-darwin-amd64 .
	cd ./dist/ && cp ./tfopen-darwin-amd64 ./tfopen && tar cfz tfopen-darwin-amd64.tar.gz ./tfopen

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build \
			 -ldflags "-X main.version=$(BUILD_VERSION)" \
			 -o dist/tfopen-darwin-arm64 .
	cd ./dist/ && cp ./tfopen-darwin-arm64 ./tfopen && tar cfz tfopen-darwin-arm64.tar.gz ./tfopen

all: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64

lint:
	golangci-lint run

modernize:
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...

local: test lint modernize

clean:
	rm ./dist/tfopen-*

