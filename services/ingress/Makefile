.PHONY: deploy-k8s finish-k8s

build-image:
	@echo "Nothing to do"

deploy-k8s:
	kubectl create -f k8s/

finish-k8s:
	kubectl delete ing main-ingress
	kubectl delete deployment default-http-backend main-ingress-lb
	kubectl delete service default-http-backend main-ingress
