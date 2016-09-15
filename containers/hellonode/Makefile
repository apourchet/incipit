MICROSERVICE="hellonode"

CONFIG='../../gcloud.json'
ACCOUNT=$(shell jq -r '.account' $(CONFIG))
PROJECT=$(shell jq -r '.project' $(CONFIG))

DOCKER_CONTAINER=$(MICROSERVICE)
DOCKER_IMGNAME="gcr.io/$(PROJECT)/$(MICROSERVICE):v1"

KUBE_CONFIG_TOOL="../../tools/kube-config.go"
KUBE_CONFIG='../../kubeconfigs/local.json'

docker-build:
	docker build -t $(DOCKER_IMGNAME) .

docker-publish:
	gcloud --account $(ACCOUNT) docker push $(DOCKER_IMGNAME) 
