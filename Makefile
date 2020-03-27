NAMESPACE ?= kubesphere-logging-system

REGISTRY?=kubespheredev
REPO_OPERATOR?=$(REGISTRY)/kube-events-operator
REPO_EXPORTER?=$(REGISTRY)/kube-events-exporter
REPO_RULER?=$(REGISTRY)/kube-events-ruler
TAG?=v0.1

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


deploy: yaml
	kubectl apply -f deploy/operator.yaml
	kubectl apply -f deploy/crs-rules-default.yaml
	kubectl apply -f deploy/crs.yaml

yaml: manifests ca-secret update-cert
	cd config/manager && $(GOBIN)/kustomize edit set image operator=$(REPO_OPERATOR):$(TAG)
	$(GOBIN)/kustomize build config/default > deploy/operator.yaml

test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=events-operator webhook paths="./pkg/..." output:crd:artifacts:config=config/crd/bases

fmt:
	go fmt ./...

vet:
	go vet ./...

generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./pkg/..."

controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

image-push: image
	docker push $(REPO_OPERATOR):$(TAG)
	docker push $(REPO_EXPORTER):$(TAG)
	docker push $(REPO_RULER):$(TAG)

.PHONY: image
image: operator-image exporter-image ruler-image

operator-image: cmd/operator/Dockerfile
	docker build -t $(REPO_OPERATOR):$(TAG) -f cmd/operator/Dockerfile .

exporter-image: cmd/exporter/Dockerfile
	docker build -t $(REPO_EXPORTER):$(TAG) -f cmd/exporter/Dockerfile .

ruler-image: cmd/ruler/Dockerfile
	docker build -t $(REPO_RULER):$(TAG) -f cmd/ruler/Dockerfile .


ca-secret:
	./hack/certs.sh --service events-webhook-service --namespace $(NAMESPACE)

update-cert: ca-secret
	./hack/update-cert.sh
