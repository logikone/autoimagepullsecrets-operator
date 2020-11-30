
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

# Install CRDs into a cluster
install: schemapatch
	kubectl apply -f deploy/crds

# Uninstall CRDs from a cluster
uninstall: schemapatch
	kubectl delete -f deploy/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

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

update-chart: rbac schemapatch

crds: crdbases
	./hack/build-crds.sh
	$(CONTROLLER_GEN) \
		schemapatch:manifests="./deploy/chart/autoimagepullsecrets-operator/charts/crds/templates" \
		output:schemapatch:artifacts:config="./deploy/chart/autoimagepullsecrets-operator/charts/crds/templates" \
		paths="./api/..."

crdbases:
	$(CONTROLLER_GEN) $(CRD_OPTIONS) paths="./api/..." output:crd:artifacts:config=./deploy/crd/bases

rbac:
	$(CONTROLLER_GEN) 'rbac:roleName=`{{ .Release.Name }}-aips-operator`' paths="./..." output:rbac:artifacts:config=./deploy/chart/autoimagepullsecrets-operator/templates

schemapatch:
	$(CONTROLLER_GEN) schemapatch:manifests=./deploy/crds output:schemapatch:artifacts:config=./deploy/crds paths=./api/...