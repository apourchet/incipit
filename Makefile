# Top level Makefile
# Handles cluster creation and docker vm creation

PROJECT_NAME=$(shell jq -r '.name' project.json)

# GCLOUD THINGS
ACCOUNT=$(shell jq -r '.account' gcloud.json)
PROJECT=$(shell jq -r '.project' gcloud.json)
ZONE=$(shell jq -r '.zone' gcloud.json)
GCLOUD_OPTS=--account $(ACCOUNT) --project $(PROJECT)

# CLUSTER THINGS
CLUSTER_NAME=$(PROJECT_NAME)
CLUSTER_NODES=$(shell jq -r '.cluster.nodes' project.json)
DOWNCMD="mount | grep -o 'on /var/lib/kubelet.* type' | cut -c 4- | rev | cut -c 6- | rev | sort -r  | xargs --no-run-if-empty sudo umount"

# DOCKER THINGS
DOCKER_MACHINE_IP=$(shell docker-machine ip dummy)

default:
	@echo $(PROJECT_NAME)
	@echo $(ACCOUNT)
	@echo $(PROJECT)
	@echo $(ZONE)
	@echo $(CLUSTER_NAME)
	@echo $(CLUSTER_NODES)
	@echo $(DOCKER_MACHINE_IP)

build-images:
	make -C services/hellogo build-image
	make -C services/ingress build-image

gcloud-kup:
	gcloud $(GCLOUD_OPTS) config set compute/zone $(ZONE)
	gcloud $(GCLOUD_OPTS) container clusters create $(CLUSTER_NAME) --num-nodes $(CLUSTER_NODES)
	gcloud $(GCLOUD_OPTS) container clusters get-credentials $(CLUSTER_NAME)

gcloud-kdown:
	gcloud $(GCLOUD_OPTS) container clusters delete $(CLUSTER_NAME)

docker-create-vm:
	docker-machine create -d virtualbox $(PROJECT_NAME);
	@echo "Now execute in shell:\neval (docker-machine env $(PROJECT_NAME))"

kup:
	kubectl config set-cluster $(CLUSTER_NAME) --server=http://$(DOCKER_MACHINE_IP):8080
	kubectl config set-context $(PROJECT_NAME) --cluster=$(CLUSTER_NAME) --namespace=$(PROJECT_NAME)
	kubectl config use-context $(PROJECT_NAME)
	docker-compose -f kubemaster/docker-compose.yaml up -d
	sleep 40
	kubectl create namespace $(PROJECT_NAME)
	make -C kubedns deploy-k8s
	@echo "--------------------------------"
	@echo "kubectl cluster-info"

kdown:
	docker-compose -f kubemaster/docker-compose.yaml down
	docker-machine ssh $(PROJECT_NAME) $(DOWNCMD)
	docker ps -a -f "name=k8s_" -q | xargs docker rm -f

deploy-k8s:
	make -C services/ingress deploy-k8s
	make -C services/hellogo deploy-k8s

finish-k8s:
	make -C services/ingress finish-k8s
	make -C services/hellogo finish-k8s


