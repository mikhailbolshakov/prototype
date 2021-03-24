#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

# install chart
helm repo add bitnami https://charts.bitnami.com/bitnami

kubectl apply -f ./cm-sqlinitscripts.yaml

helm install -f ./values.yaml pg bitnami/postgresql