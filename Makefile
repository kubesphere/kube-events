NAMESPACE ?= kubesphere-logging-system

REGISTRY?=kubesphere
REPO_OPERATOR?=$(REGISTRY)/kube-events-operator
REPO_EXPORTER?=$(REGISTRY)/kube-events-exporter
REPO_RULER?=$(REGISTRY)/kube-events-ruler
TAG?=latest

GO_PKG?=github.com/kubesphere/kube-events

CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

CONTROLLER_GEN := $(GOBIN)/controller-gen
KE_DOCGEN_BINARY:=$(GOBIN)/kube-events-docgen

TYPES_V1ALPHA1_TARGET := pkg/apis/v1alpha1/exporter_types.go pkg/apis/v1alpha1/ruler_types.go pkg/apis/v1alpha1/rule_types.go
DEEPCOPY_TARGET := pkg/apis/v1alpha1/zz_generated.deepcopy.go

.PHONY: helm

deploy:
	kubectl apply -f config/bundle.yaml

generate: $(DEEPCOPY_TARGET) manifests
	cd config && $(GOBIN)/kustomize edit set image operator=$(REPO_OPERATOR):$(TAG)
	$(GOBIN)/kustomize build config > config/bundle.yaml
	cd config/crs && $(GOBIN)/kustomize edit set image exporter=$(REPO_EXPORTER):$(TAG) ruler=$(REPO_RULER):$(TAG)
	$(GOBIN)/kustomize build config/crs > config/crs/bundle.yaml

# Generate manifests e.g. CRD, RBAC etc.
manifests: $(CONTROLLER_GEN) $(TYPES_V1ALPHA1_TARGET)
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=operator paths="./pkg/..." output:crd:artifacts:config=config/crd/bases

helm: $(CONTROLLER_GEN)
	kustomize build config/helm | sed -e '/creationTimestamp/d' > helm/crds/bundle.yaml
	tar zcvf kube-events.tgz helm

doc/api.md: $(KE_DOCGEN_BINARY) $(TYPES_V1ALPHA1_TARGET)
	$(KE_DOCGEN_BINARY) $(TYPES_V1ALPHA1_TARGET) > doc/api.md

fmt:
	go fmt ./...

vet:
	go vet ./...

# Generate deepcopy etc.
$(DEEPCOPY_TARGET): $(CONTROLLER_GEN) $(TYPES_V1ALPHA1_TARGET)
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./pkg/..."

$(CONTROLLER_GEN):
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\

$(KE_DOCGEN_BINARY): cmd/docgen/kube-events-docgen.go
	go install cmd/docgen/kube-events-docgen.go

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

cross-build: cross-build-operator cross-build-exporter cross-build-ruler

cross-build-operator: cmd/operator/Dockerfile
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(REPO_OPERATOR):$(TAG) -f cmd/operator/Dockerfile .

cross-build-exporter: cmd/exporter/Dockerfile
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(REPO_EXPORTER):$(TAG) -f cmd/exporter/Dockerfile .

cross-build-ruler: cmd/ruler/Dockerfile
	docker buildx build --push --platform linux/amd64,linux/arm64 -t $(REPO_RULER):$(TAG) -f cmd/ruler/Dockerfile .

ca-secret:
	./hack/certs.sh --service kube-events-admission --namespace $(NAMESPACE)

update-cert: ca-secret
	./hack/update-cert.sh
