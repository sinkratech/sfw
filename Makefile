build_macOS_x86:
	GOOS=darwin GOARCH=amd64 go build -o sfw-macOS-amd64 cmd/sfw/main.go

build_macOS_arm:
	GOOS=darwin GOARCH=arm64 go build -o sfw-macOS-arm64 cmd/sfw/main.go

build_macOS: build_macOS_arm build_macOS_x86

build_linux_x86:
	GOOS=linux GOARCH=amd64 go build -o sfw-linux-amd64 cmd/sfw/main.go

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build  -o sfw-linux-arm64 cmd/sfw/main.go

build_linux_386:
	GOOS=linux GOARCH=386 go build -o sfw-linux-386 cmd/sfw/main.go

build_linux: build_linux_386 build_linux_arm64 build_linux_x86 

build: build_linux build_macOS

tar:
	tar zcvf sfw.tar.gz sfw-macOS-amd64 sfw-macOS-arm64 sfw-linux-amd64 sfw-linux-arm64 sfw-linux-386