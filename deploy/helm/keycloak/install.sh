#!/bin/bash

helm repo add bitnami https://charts.bitnami.com/bitnami
kubectl apply -f ./keycloak-sqlinitscripts.yaml
helm install -f ./values.yaml keycloak bitnami/keycloak