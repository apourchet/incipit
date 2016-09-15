MICROSERVICE="hellogo"

CONFIG='../../gcloud.json'
ACCOUNT=$(shell jq -r '.account' $(CONFIG))
PROJECT=$(shell jq -r '.project' $(CONFIG))

DOCKER_CONTAINER=$(MICROSERVICE)
DOCKER_IMGNAME="gcr.io/$(PROJECT)/$(MICROSERVICE):v1"

KUBE_CONFIG_TOOL="../../tools/kube-config.go"
KUBE_CONFIG='../../kubeconfigs/local.json'

default:
	go build -o bin/hellogo server.go 

build-image:
	docker build -t $(DOCKER_IMGNAME) .

docker-push-img:
	gcloud --account $(ACCOUNT) docker push $(DOCKER_IMGNAME) 

k8s-resources:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*-res.yaml | kubectl create -f -

k8s-services:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*-svc.yaml | kubectl create -f -

k8s-deployments:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*-dp.yaml | kubectl create -f -

k8s-undeploy:
	kubectl delete deployment $(MICROSERVICE)

none:
	@echo $(ACCOUNT)
	@echo $(PROJECT)
	@echo $(DOCKER_IMGNAME)
