.PHONY: docker-build

CONFIG = ../gcloud.json
ACCOUNT = $(shell jq -r '.account' $(CONFIG))
PROJECT = $(shell jq -r '.project' $(CONFIG))

DOCKER_IMAGE_PREFIX = gcr.io/$(PROJECT)
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
	docker build -f $*/Dockerfile -t $(DOCKER_IMAGE_PREFIX)/$*:v1 ..

# Build all docker images for all containers
dockerize:
	$(foreach container, $(CONTAINERS), \
	docker build -f $(container)/Dockerfile -t $(DOCKER_IMAGE_PREFIX)/$(container):v1 ..;)
