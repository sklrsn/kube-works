.DEFAULT_GOAL:all

.PHONY: all
all: down build up

.PHONY: build
build:
	docker compose -f docker-compose.yaml build

.PHONY: up
up:
	docker compose -f docker-compose.yaml up

.PHONY: down
down:
	docker compose -f docker-compose.yaml down

.PHONY: tag
tag:
	docker tag kube-two-backend:latest sklrsn/kube-backend:latest
	docker tag kube-two-frontend:latest sklrsn/kube-frontend:latest
	docker tag kube-two-postgres-server:latest sklrsn/kube-postgres:latest
	docker tag kube-two-redis-server:latest sklrsn/kube-redis:latest  

.PHONY: push
push: build tag
	docker push sklrsn/kube-backend:latest
	docker push sklrsn/kube-frontend:latest
	docker push sklrsn/kube-postgres:latest
	docker push sklrsn/kube-redis:latest  

.PHONY: deploy
deploy:
	kubectl apply -f k8s

.PHONY: delete
delete:
	kubectl delete -f k8s

.PHONY: ingress-controller
ingress-controller:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.6.4/deploy/static/provider/cloud/deploy.yaml