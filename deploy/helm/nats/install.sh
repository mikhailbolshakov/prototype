#!/bin/bash

helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm install -f ./values.yaml nats nats/stan