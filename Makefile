# Top level Makefile
# Handles cluster creation and docker vm creation

PROJECT_NAME  =  $(shell jq -r '.name' project.json)

# CLUSTER THINGS
CLUSTER_NAME = $(PROJECT_NAME)
DOWNCMD = "mount | grep -o 'on /var/lib/kubelet.* type' | cut -c 4- | rev | cut -c 6- | rev | sort -r  | xargs --no-run-if-empty sudo umount"

# DOCKER THINGS
DOCKER_MACHINE_NAME = $(PROJECT_NAME)
DOCKER_MACHINE_IP = $(shell docker-machine ip $(DOCKER_MACHINE_NAME))
DOCKER_BUILDER_IMAGENAME = apourchet/golang-grpc
DOCKER_BUILDER_CONTAINER = $(PROJECT_NAME)-builder-container
ETC_HOST_HACK_UNDO = sudo sed -i '' "/$(DOCKER_MACHINE_NAME)\.machine/d" /etc/hosts
ETC_HOST_HACK_DO = 

# PROTOC THINGS
PROTOC_OPTS = -I protos
PROTOC_OPTS += -I /usr/local/include
PROTOC_OPTS += -I $(GOPATH)/src
PROTOC_OPTS += -I $(shell pwd)/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
PROTOC_OPTS += --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:protos/go/

# TESTER
TESTER_NAME = tester

# TOOLS
KUBE_CONFIG_TOOL = ./tools/kubeconfig/kubeconfig.go
KUBE_CONFIG = ./kubeconfigs/local.json
KUBERNETES_PUBLIC_ADDRESS=$(gcloud compute addresses describe kubernetes --format 'value(address)')

.PHONY: resources deployments docker-builder protos

# -------BUILDING
# Create docker host virtual machine
docker-create-vm:
	docker-machine create -d virtualbox $(DOCKER_MACHINE_NAME)
	$(ETC_HOST_HACK_UNDO)
	sudo /bin/bash -c "echo \"`docker-machine ip $(DOCKER_MACHINE_NAME)`    $(DOCKER_MACHINE_NAME).machine\" >> /etc/hosts"
	docker-machine scp /usr/local/docker-dev/machine/bootsync.sh $(DOCKER_MACHINE_NAME):/tmp/bootsync.sh
	docker-machine ssh $(DOCKER_MACHINE_NAME) sudo mv /tmp/bootsync.sh /var/lib/boot2docker/bootsync.sh
	docker-machine ssh $(DOCKER_MACHINE_NAME) /var/lib/boot2docker/bootsync.sh
	@echo "---------------------------------------------------------"
	@echo "Now execute in shell:\neval (docker-machine env $(DOCKER_MACHINE_NAME))"

# Destroy that VM
docker-destroy-vm:
	docker-machine rm $(DOCKER_MACHINE_NAME)

# Command that should only be run from within a container
# The following are general build targets
build:
	make protos
	make -C containers build

protos:
	ldconfig
	which protoc || make -C ./vendor/github.com/google/protobuf install
	protoc $(PROTOC_OPTS) protos/*.proto 
	protoc $(PROTOC_OPTS) --grpc-gateway_out=logtostderr=true:protos/go/ protos/*.proto 
	protoc $(PROTOC_OPTS) --swagger_out=logtostderr=true:protos/swagger/ protos/*.proto 

# Create a docker container that will be used 
# to do all of the building
# docker build -f Dockerfile -t $(DOCKER_BUILDER_IMAGENAME) .
docker-builder:
	-docker rm $(DOCKER_BUILDER_CONTAINER) -f
	docker run --name $(DOCKER_BUILDER_CONTAINER) \
		-v /Users/antoine/gopath/src/github.com/apourchet/incipit:/go/src/github.com/apourchet/incipit \
		-d \
		$(DOCKER_BUILDER_IMAGENAME) /bin/sh -c "while true; do sleep 10; done"

# Build all within docker
# Then dockerize all of the containers
docker-build:
	docker exec $(DOCKER_BUILDER_CONTAINER) make build
	make dockerize

docker-protos:
	docker exec $(DOCKER_BUILDER_CONTAINER) make protos

dockerize:
	make -C containers dockerize

# -------KUBERNETES
# Spin up the k8s cluster in the docker VM
kup:
	kubectl config set-cluster $(CLUSTER_NAME) --server http://$(DOCKER_MACHINE_IP):8080
	kubectl config set-context $(PROJECT_NAME) --cluster $(CLUSTER_NAME) --namespace $(PROJECT_NAME)
	kubectl config use-context $(PROJECT_NAME)
	docker-compose -f kubemaster/docker-compose.yaml up -d
	bash ./tools/retry.sh "kubectl cluster-info" 2
	-kubectl create namespace $(PROJECT_NAME)

# Tear down the k8s cluster in the docker VM
kdown:
	docker-compose -f kubemaster/docker-compose.yaml down
	docker-machine ssh $(DOCKER_MACHINE_NAME) "sudo rm -rf /etcd-data"
	docker-machine ssh $(DOCKER_MACHINE_NAME) $(DOWNCMD)
	docker ps -a -f "name = k8s_" -q | xargs docker rm -f

# Create services, secrets and persistent disks
resources:
	-kubectl create namespace $(PROJECT_NAME)
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./resources/*/*.json | kubectl apply -f -

# Create deployments/replication controllers/pods
deployments:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./deployments/*/*.json | kubectl apply -f -

# Delete all deployments
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

# Start the Kubernetes Dashboard
ui-k8s:
	kubectl create -f https://rawgit.com/kubernetes/dashboard/master/src/deploy/kubernetes-dashboard.yaml

# Create local certificates for TLS
local-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./misc/local-server.key -out ./misc/local-server.crt -subj "/CN=$(DOCKER_MACHINE_NAME).machine"

stage-certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./misc/stage-server.key -out ./misc/stage-server.crt -subj "/CN=$(KUBERNETES_PUBLIC_ADDRESS)"

# -------SERVICE SPECIFIC TARGETS
# Should be run within the docker builder container
build-%:
	make -C containers build-$*

# Builds and dockerizes the target
docker-build-%:
	docker exec $(DOCKER_BUILDER_CONTAINER) make build-$*
	make -C containers dockerize-$*

# Deletes the designated pod
recall-%:
	kubectl delete deployment $* 
	make deployments

# Deploys the resources and deployments for a service
deploy-%:
	-go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./resources/$*/*.json | kubectl apply -f -
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./deployments/$*/*.json | kubectl apply -f -

# Deletes a pod and bounces the ingress pod as well
bounce-%:
	kubectl get pods | grep $* | cut -f 1 -d ' ' | tail -n +1 | xargs kubectl delete pod 
	kubectl get pods | grep ingress | cut -f 1 -d ' ' | tail -n 1 | xargs kubectl delete pod 
	sleep 1

# Rebuild and redeploy a service
loop-%: 
	make docker-build-$*
	make bounce-$*

