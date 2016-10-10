# Top level Makefile
# Handles cluster creation and docker vm creation

PROJECT_NAME  =  $(shell jq -r '.name' project.json)

# GCLOUD THINGS
ACCOUNT = $(shell jq -r '.account' gcloud.json)
PROJECT = $(shell jq -r '.project' gcloud.json)
ZONE = $(shell jq -r '.zone' gcloud.json)
GCLOUD_OPTS = --account $(ACCOUNT) --project $(PROJECT)

# CLUSTER THINGS
CLUSTER_NAME = $(PROJECT_NAME)
CLUSTER_NODES = $(shell jq -r '.cluster.nodes' project.json)
DOWNCMD = "mount | grep -o 'on /var/lib/kubelet.* type' | cut -c 4- | rev | cut -c 6- | rev | sort -r  | xargs --no-run-if-empty sudo umount"

# DOCKER THINGS
DOCKER_MACHINE_NAME = $(PROJECT_NAME)
DOCKER_MACHINE_IP = $(shell docker-machine ip $(DOCKER_MACHINE_NAME))
DOCKER_BUILDER_IMAGENAME = $(PROJECT_NAME)-builder-image
DOCKER_BUILDER_CONTAINER = $(PROJECT_NAME)-builder-container
ETC_HOST_HACK_UNDO = sudo sed -i '' "/$(DOCKER_MACHINE_NAME)\.machine/d" /etc/hosts
ETC_HOST_HACK_DO = 

# TESTER
TESTER_NAME = tester

# TOOLS
KUBE_CONFIG_TOOL = ./tools/kubeconfig/kubeconfig.go
KUBE_CONFIG = ./kubeconfigs/local.json

.PHONY: resources deployments docker-builder

docker-create-vm:
	docker-machine create -d virtualbox $(DOCKER_MACHINE_NAME)
	$(ETC_HOST_HACK_UNDO)
	sudo /bin/bash -c "echo \"`docker-machine ip $(DOCKER_MACHINE_NAME)`    $(DOCKER_MACHINE_NAME).machine\" >> /etc/hosts"
	docker-machine scp /usr/local/docker-dev/machine/bootsync.sh $(DOCKER_MACHINE_NAME):/tmp/bootsync.sh
	docker-machine ssh $(DOCKER_MACHINE_NAME) sudo mv /tmp/bootsync.sh /var/lib/boot2docker/bootsync.sh
	docker-machine ssh $(DOCKER_MACHINE_NAME) /var/lib/boot2docker/bootsync.sh
	@echo "---------------------------------------------------------"
	@echo "Now execute in shell:\neval (docker-machine env $(DOCKER_MACHINE_NAME))"

docker-destroy-vm:
	docker-machine rm $(DOCKER_MACHINE_NAME)

# Command that should only be run from within a container
# The following are general build targets
build:
	make -C containers build

docker-builder:
	docker build -f Dockerfile -t $(DOCKER_BUILDER_IMAGENAME) .
	-docker rm $(DOCKER_BUILDER_CONTAINER) -f
	docker run --name $(DOCKER_BUILDER_CONTAINER) \
		-v /Users/antoine/gopath/src/github.com/apourchet/incipit:/go/src/github.com/apourchet/incipit \
		-d \
		$(DOCKER_BUILDER_IMAGENAME) /bin/sh -c "while true; do sleep 10; done"

docker-build:
	docker exec $(DOCKER_BUILDER_CONTAINER) make build
	make dockerize

dockerize:
	make -C containers dockerize

kup:
	kubectl config set-cluster $(CLUSTER_NAME) --server http://$(DOCKER_MACHINE_IP):8080
	kubectl config set-context $(PROJECT_NAME) --cluster $(CLUSTER_NAME) --namespace $(PROJECT_NAME)
	kubectl config use-context $(PROJECT_NAME)
	docker-compose -f kubemaster/docker-compose.yaml up -d
	bash ./tools/retry.sh "kubectl cluster-info" 2
	-kubectl create namespace $(PROJECT_NAME)

kdown:
	docker-compose -f kubemaster/docker-compose.yaml down
	docker-machine ssh $(DOCKER_MACHINE_NAME) "sudo rm -rf /etcd-data"
	docker-machine ssh $(DOCKER_MACHINE_NAME) $(DOWNCMD)
	docker ps -a -f "name = k8s_" -q | xargs docker rm -f

resources:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./resources/*/*.json | kubectl apply -f -

deployments:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./deployments/*/*.json | kubectl apply -f -

recall:
	kubectl get deployments | cut -f 1 -d ' ' | tail -n +2 | xargs kubectl delete deployments

# Spin up pod that tests
test: docker-build-tester
	-kubectl delete job tester-job
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./jobs/tester/*.json | kubectl apply -f -

# Test locally
local-test:
	go test -v ./lib/...
	go test -v ./containers/...

ui:
	kubectl create -f https://rawgit.com/kubernetes/dashboard/master/src/deploy/kubernetes-dashboard.yaml

local-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./misc/local-server.key -out ./misc/local-server.crt -subj "/CN=$(DOCKER_MACHINE_NAME).machine"

# Service specific targets
build-%:
	make -C containers build-$*

docker-build-%:
	docker exec $(DOCKER_BUILDER_CONTAINER) make build-$*
	make -C containers dockerize-$*

recall-%:
	kubectl get pods | grep $* | cut -f 1 -d ' ' | tail -n 1 | xargs kubectl delete pod 

deploy-%:
	-go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./resources/$*/*.json | kubectl apply -f -
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./deployments/$*/*.json | kubectl apply -f -

bounce-%:
	make recall-$*
	kubectl get pods | grep ingress | cut -f 1 -d ' ' | tail -n 1 | xargs kubectl delete pod 
	sleep 1

loop-%: 
	make docker-build-$*
	make bounce-$*

# GOOGLE SPECIFIC TARGETS
gcloud-kup:
	gcloud $(GCLOUD_OPTS) config set compute/zone $(ZONE)
	gcloud $(GCLOUD_OPTS) container clusters create $(CLUSTER_NAME) --num-nodes $(CLUSTER_NODES)
	gcloud $(GCLOUD_OPTS) container clusters get-credentials $(CLUSTER_NAME)

gcloud-kdown:
	gcloud $(GCLOUD_OPTS) container clusters delete $(CLUSTER_NAME)


