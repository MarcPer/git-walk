.PHONY: build

build:
	GOARCH="amd64" GOOS="linux" go build -o ./build/linux/git-walk
	GOARCH="amd64" GOOS="darwin" go build -o ./build/darwin/git-walk
