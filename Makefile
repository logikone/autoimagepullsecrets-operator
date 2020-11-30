
# Image URL to use all building/pushing image targets
IMG ?= logikone/autoimagepullsecrets-operator:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:crdVersions=v1beta1"

CONTROLLER_GEN ?= go run sigs.k8s.io/controller-tools/cmd/controller-gen

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: operator

clean:
	find . -type f -name c.out -delete

# Run tests
test: generate fmt vet manifests
	go run github.com/onsi/ginkgo/ginkgo \
		-v \
		-cover \
		-coverprofile c.out \
		-outputdir . \
		./...

# Build manager binary
operator: generate fmt vet
	go build -o bin/operator ./cmd/operator

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./cmd/operator

# Generate manifests e.g. CRD, RBAC etc.
manifests: update-chart
	#$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate:
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

update-chart: crds deploy/rbac/role.yaml deploy/chart/autoimagepullsecrets-operator/templates/role.yaml

crds:
	$(CONTROLLER_GEN) $(CRD_OPTIONS) paths="./api/..." output:crd:artifacts:config=./deploy/crd/bases
	@./hack/build-crds.sh

.PHONY: deploy/rbac/role.yaml
deploy/rbac/role.yaml:
	$(CONTROLLER_GEN) 'rbac:roleName=aips-operator' paths="./..." output:rbac:artifacts:config=./deploy/rbac/

.PHONY: deploy/webhooks/manifests.yaml
deploy/webhooks/manifests.yaml:
	$(CONTROLLER_GEN) webhook paths="./..." output:webhook:artifacts:config=./deploy/webhooks/

deploy/chart/autoimagepullsecrets-operator/templates/webhooks.yaml: deploy/webhooks/manifests.yaml deploy/webhooks/yq-scripts.yaml
	yq merge -x -d '*' deploy/webhooks/manifests.yaml deploy/chart/labels.yaml | \
		yq write -d '*' -s deploy/webhooks/yq-scripts.yaml -- /dev/stdin > $@

deploy/chart/autoimagepullsecrets-operator/templates/role.yaml: deploy/rbac/role.yaml
	yq merge -x -d '*' deploy/rbac/role.yaml deploy/chart/labels.yaml | \
		yq write -d '*' -s deploy/rbac/yq-scripts.yaml -- /dev/stdin > $@