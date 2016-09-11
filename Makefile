# Top level Makefile
# Handles cluster creation and docker vm creation
# 

PROJECT_NAME="dummy"

# GCLOUD THINGS
ACCOUNT=$(shell jq -r '.gcloud_account' config.json)
PROJECT=$(shell jq -r '.gcloud_project' config.json)
ZONE=$(shell jq -r '.gcloud_zone' config.json)
GCLOUD_OPTS=--account $(ACCOUNT) --project $(PROJECT)

# CLUSTER THINGS
CLUSTER_NAME=$(PROJECT_NAME)
CLUSTER_NODES=2
DOWNCMD="mount | grep -o 'on /var/lib/kubelet.* type' | cut -c 4- | rev | cut -c 6- | rev | sort -r  | xargs --no-run-if-empty sudo umount"

# DOCKER THINGS
DOCKER_MACHINE_IP=$(shell docker-machine ip dummy)

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
	@kubectl config set-cluster $(CLUSTER_NAME) --server=http://$(DOCKER_MACHINE_IP):8080
	@kubectl config set-context $(CLUSTER_NAME) --cluster=$(CLUSTER_NAME)
	@docker-compose -f kubemaster/docker-compose.yaml up -d
	@echo "--------------------------------"
	@echo "kubectl config use-context $(CLUSTER_NAME)"
	@echo "kubectl cluster-info"

kdown:
	@docker-compose -f kubemaster/docker-compose.yaml down
	@docker-machine ssh $(PROJECT_NAME) $(DOWNCMD)
	@docker ps -a -f "name=k8s_" -q | xargs docker rm -f

none:
	@echo $(ACCOUNT)
	@echo $(PROJECT)
	@echo $(CLUSTER_NAME)
	@echo $(ZONE)

