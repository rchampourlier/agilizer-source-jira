GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
CMD_NAME=agilizer-source-jira
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/$(CMD_NAME)
DOCKER_TAG=latest

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

container: $(DOCKER_CMD)
	docker build . -t $(CMD_NAME):$(DOCKER_TAG)

clean:
	rm -rf $(DOCKER_BUILD)