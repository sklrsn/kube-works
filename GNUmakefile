.DEFAULT_GOAL: all

.PHONY: all
all: build push

.PHONY: build
build:
	docker build -t sklrsn/kube-welcome:latest .

.PHONY: push
push:
	docker push sklrsn/kube-welcome:latest

.PHONY: deploy
deploy:
	kubectl apply -f kube-welcome-pod.yaml
	kubectl apply -f kube-welcome-service.yaml

.PHONY: delete
delete:
	kubectl delete -f kube-welcome-pod.yaml
	kubectl delete -f kube-welcome-service.yaml