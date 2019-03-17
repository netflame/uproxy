APPNAME := uproxy

.PHONY: build
build: build_darwin

.PHONY: build_darwin
build_darwin: test clean
	@GOOS=darwin GOARCH=amd64 \
		sh build.sh

.PHONY: build_linux
build_linux: test clean
	@GOOS=linux GOARCH=amd64 \
		sh build.sh

.PHONY: test
test:
	@go test

# for docker
.PHONY: docker_run
docker_run: docker_build
	@docker-compose up -d

.PHONY: docker_build
docker_build: build_linux
	@docker-compose build && make clean

.PHONY: clean
clean:
	@go mod tidy && go clean && rm -f ${APPNAME}*
