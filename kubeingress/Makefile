
MICROSERVICE="kubeingress"

CONFIG='../config.json'
GCLOUD_ACCOUNT=$(shell jq -r '.gcloud_account' $(CONFIG))
GCLOUD_PROJECT=$(shell jq -r '.gcloud_project' $(CONFIG))

DOCKER_CONTAINER='kubeingress'
DOCKER_IMGNAME="gcr.io/$(GCLOUD_PROJECT)/kubeingress:v1"

build-image:
	docker build -t $(DOCKER_IMGNAME) .

docker-push-img:
	gcloud --account $(GCLOUD_ACCOUNT) docker push $(DOCKER_IMGNAME) 

deploy-k8s:
	kubectl create -f .

finish-k8s:
	kubectl delete ing main-ingress
	kubectl delete deployment default-http-backend main-ingress-lb
	kubectl delete service default-http-backend main-ingress

none:
	@echo $(GCLOUD_ACCOUNT)
	@echo $(GCLOUD_PROJECT)
	@echo $(DOCKER_IMGNAME)
