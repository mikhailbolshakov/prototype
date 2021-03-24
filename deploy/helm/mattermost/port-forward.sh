#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

kubectl port-forward --namespace default $(kubectl get pods --namespace default -l "app=mattermost" -o jsonpath='{ .items[0].metadata.name }') 8065:8065
