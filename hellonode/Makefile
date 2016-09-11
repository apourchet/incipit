
CONFIG='../config.json'
GCLOUD_ACCOUNT=$(shell jq -r '.gcloud_account' $(CONFIG))
GCLOUD_PROJECT=$(shell jq -r '.gcloud_project' $(CONFIG))
CLUSTER_NAME=$(shell jq -r '.k8s_cluster_name' $(CONFIG))

DOCKER_CONTAINER='hellonode'
DOCKER_IMGNAME="gcr.io/$(GCLOUD_PROJECT)/hellonode:v1"

build-image:
	docker build -t $(DOCKER_IMGNAME) .

run-local:
	docker run -d -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMGNAME)

stop-local:
	docker stop $(DOCKER_CONTAINER)

deploy-k8s:
	kubectl create -f ./k8s

finish-k8s:
	kubectl delete service hellonode-svc
	kubectl delete deployment hellonode-dp

docker-push-img:
	gcloud --account $(GCLOUD_ACCOUNT) docker push $(DOCKER_IMGNAME) 

none:
	@echo $(GCLOUD_ACCOUNT)
	@echo $(GCLOUD_PROJECT)
	@echo $(CLUSTER_NAME)
	@echo $(DOCKER_IMGNAME)
