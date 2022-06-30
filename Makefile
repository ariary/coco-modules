before.build:
	go mod tidy && go mod download

build.module.hello:
	@echo "build in ${PWD}";go build cmd/hello/hello.go