.PHONY: deploy-k8s finish-k8s

deploy-k8s:
	kubectl create -f .

finish-k8s:
	kubectl delete ing main-ingress
	kubectl delete deployment default-http-backend main-ingress-lb
	kubectl delete service default-http-backend main-ingress
