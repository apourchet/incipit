.PHONY: docker-build

CONFIG = ../project.json
PROJECT = $(shell jq -r '.name' $(CONFIG))

DOCKER_IMAGE_PREFIX = 
CONTAINERS = auth tester logger

.PHONY: build

# Build the binaries and whatnot on the current machine
# for the specified container
build-%:
	cd $* && sh build.sh && cd -

# Build all binaries for all containers
build:
	$(foreach container, $(CONTAINERS), \
	cd $(container) && sh build.sh && cd -;)

# Build docker image for the specified container
dockerize-%:
	docker build -f $*/Dockerfile -t $(DOCKER_IMAGE_PREFIX)$*:v1 ..
	docker tag $(DOCKER_IMAGE_PREFIX)$*:v1 $(DOCKER_IMAGE_PREFIX)$*:latest

# Build all docker images for all containers
dockerize:
	$(foreach container, $(CONTAINERS), \
	docker build -f $(container)/Dockerfile -t $(DOCKER_IMAGE_PREFIX)$(container):v1 ..; \
	docker tag $(DOCKER_IMAGE_PREFIX)$(container):v1 $(DOCKER_IMAGE_PREFIX)$(container):latest;)
