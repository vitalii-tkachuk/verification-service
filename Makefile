POSTGRES = $(shell kubectl get pods -o=jsonpath='{.items[?(@.spec.containers[*].name=="postgres")].metadata.name}')
APP = $(shell kubectl get pods -o=jsonpath='{.items[?(@.spec.containers[*].name=="application")].metadata.name}')
IP = $(shell minikube ip)
PWD = $(shell pwd)

.PHONY: start stop helm-debug helm-upgrade postgres skaffold hosts format vet lint test

start:
	minikube start --vm=true --driver=hyperkit
	minikube addons enable ingress

stop:
	minikube stop

helm-debug:
	helm template --debug deployments/minikube/helm

helm-upgrade:
	helm upgrade verification-service deployments/minikube/helm --install

postgres:
	kubectl exec -it $(POSTGRES) -c postgres -- psql -U postgres verification-service

logs:
	kubectl logs $(APP) -c application -f

skaffold:
	skaffold dev

hosts:
	sudo sh -c "echo '$(IP) verification-service.local' >> /etc/hosts"

format:
	docker run --rm -v $(PWD):/app -w /app golang:1.19 go fmt ./...

vet:
	docker run --rm -v $(PWD):/app -w /app golang:1.19 go vet ./...

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run ./... -v

test:
	docker run --rm -v $(PWD):/app -w /app golang:1.19 go test -race -p 4 -vet=off ./... -v
