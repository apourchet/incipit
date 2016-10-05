.PHONY: docker-build

CONFIG = ../gcloud.json
ACCOUNT = $(shell jq -r '.account' $(CONFIG))
PROJECT = $(shell jq -r '.project' $(CONFIG))

DOCKER_IMAGE_PREFIX = gcr.io/$(PROJECT)
CONTAINERS = auth 

build:
	$(foreach container, $(CONTAINERS), \
	sh $(container)/build.sh)

build-docker:
	$(foreach container, $(CONTAINERS), \
	docker build -t $(DOCKER_IMAGE_PREFIX)/$(container):v1 $(container))
