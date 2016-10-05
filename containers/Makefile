.PHONY: docker-build

CONFIG = ../gcloud.json
ACCOUNT = $(shell jq -r '.account' $(CONFIG))
PROJECT = $(shell jq -r '.project' $(CONFIG))

DOCKER_IMAGE_PREFIX = gcr.io/$(PROJECT)
CONTAINERS = auth 

.PHONY: build

build-%:
	cd $* && sh build.sh && cd -

build:
	$(foreach container, $(CONTAINERS), \
	cd $(container) && sh build.sh && cd -)

dockerize-%:
	cd $* && docker build -t $(DOCKER_IMAGE_PREFIX)/$*:v1 . && cd -

dockerize:
	$(foreach container, $(CONTAINERS), \
	cd $(container) && docker build -t $(DOCKER_IMAGE_PREFIX)/$(container):v1 . && cd -)
