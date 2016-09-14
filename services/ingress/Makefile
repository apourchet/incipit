.PHONY: deploy-k8s finish-k8s

KUBE_CONFIG_TOOL="../../tools/kube-config.go"
KUBE_CONFIG='../../kubeconfigs/local.json'

build-image:
	@echo "Nothing to do"

deploy-k8s:
	go run $(KUBE_CONFIG_TOOL) $(KUBE_CONFIG) ./k8s/*.yaml | kubectl create -f -

finish-k8s:
	kubectl delete ing main-ingress
	kubectl delete deployment default-http-backend main-ingress-lb
	kubectl delete service default-http-backend main-ingress
