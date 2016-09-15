.PHONY: deploy-k8s finish-k8s

KUBE_CONFIG_TOOL="../../tools/kube-config.go"
KUBE_CONFIG='../../kubeconfigs/local.json'

build-image:
	@echo "Nothing to do"

k8s-resources:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*-res.yaml | kubectl create -f -

k8s-services:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*-svc.yaml | kubectl create -f -

k8s-deployments:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*-dp.yaml | kubectl create -f -

k8s-undeploy:
	kubectl delete deployment default-http-backend main-ingress-lb

