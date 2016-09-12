MICROSERVICE="hellogo"

CONFIG='../gcloud.json'
ACCOUNT=$(shell jq -r '.account' $(CONFIG))
PROJECT=$(shell jq -r '.project' $(CONFIG))

DOCKER_CONTAINER=$(MICROSERVICE)
DOCKER_IMGNAME="gcr.io/$(PROJECT)/$(MICROSERVICE):v1"

default:
	go build -o bin/hellogo server.go 

build-image:
	docker build -t $(DOCKER_IMGNAME) .

run-local:
	docker run -d -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMGNAME)

stop-local:
	docker stop $(DOCKER_CONTAINER)

deploy-k8s:
	kubectl create -f ./k8s

finish-k8s:
	kubectl delete service $(MICROSERVICE)
	kubectl delete deployment $(MICROSERVICE)

docker-push-img:
	gcloud --account $(ACCOUNT) docker push $(DOCKER_IMGNAME) 

none:
	@echo $(ACCOUNT)
	@echo $(PROJECT)
	@echo $(DOCKER_IMGNAME)
