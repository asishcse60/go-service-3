SHELL := /bin/bash

run:
	go run main.go

build:
	go build -ldflags "-X main.build=local"

#show: show docker images
show:
	docker images

#Prune: Docker all images deleted
prune:
	docker system prune

#delete: Docker image delete with id
delete:
	docker rmi f86055b3c02a

# Building containers --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \

VERSION := 1.0

all: sales

sales:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		.

# Running from within k8s/kind

KIND_CLUSTER := asish-starter-cluster

# Upgrade to latest Kind (>=v0.11): e.g. brew upgrade kind
# For full Kind v0.11 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.0
# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.22.0@sha256:b8bda84bb3a190e6e028b1760d277454a72267a5454b57db34437c34a588d047 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load-linux:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/sales-pod && kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply-linux:
	cat zarf/k8s/base/sales-pod/base-sales.yaml | kubectl apply -f -

#kind-apply: for windows yaml file info
#	kubectl apply -f zarf/k8s/base/sales-pod/base-sales.yaml

kind-apply:
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-sales:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment sales-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

#Decribe the sales pod
kind-describe:
	kubectl describe pod -l app=sales

# ==============================================================================
# Modules support
tidy:
	go mod tidy
	go mod vendor