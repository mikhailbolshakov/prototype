#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

helm repo add bitnami https://charts.bitnami.com/bitnami
helm install -f ./values.yaml redis bitnami/redis