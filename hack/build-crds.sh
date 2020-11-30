#!/usr/bin/env bash

CRD_BASES=$(find deploy/crd/bases -type f -name "*.yaml")
DESTINATION_DIR=./deploy/chart/autoimagepullsecrets-operator/charts/crds/templates

for BASE in ${CRD_BASES}; do
  yq merge "${BASE}" deploy/crd/template.yaml >"${DESTINATION_DIR}"/"$(basename "${BASE}")"
done
