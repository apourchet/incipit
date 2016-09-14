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

deploy-k8s:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*.yaml | kubectl create -f -

finish-k8s:
	kubectl delete service $(MICROSERVICE)
	kubectl delete deployment $(MICROSERVICE)

docker-push-img:
	gcloud --account $(ACCOUNT) docker push $(DOCKER_IMGNAME) 

none:
	@echo $(ACCOUNT)
	@echo $(PROJECT)
	@echo $(DOCKER_IMGNAME)
