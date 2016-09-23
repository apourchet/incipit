MICROSERVICE=hermes

CONFIG='../../gcloud.json'
ACCOUNT=$(shell jq -r '.account' $(CONFIG))
PROJECT=$(shell jq -r '.project' $(CONFIG))

DOCKER_CONTAINER=$(MICROSERVICE)
DOCKER_IMGNAME=gcr.io/$(PROJECT)/$(MICROSERVICE):v1

KUBE_CONFIG_TOOL="../../tools/kube-config.go"
KUBE_CONFIG='../../kubeconfigs/local.json'

.PHONY: build

build:
	CGO_ENABLED=0 govendor build -ldflags '-extldflags "-static"' -a -installsuffix cgo -o build/$(MICROSERVICE) server.go

build-docker:
	docker build -t $(DOCKER_IMGNAME) .

docker-publish:
	gcloud --account $(ACCOUNT) docker push $(DOCKER_IMGNAME) 
