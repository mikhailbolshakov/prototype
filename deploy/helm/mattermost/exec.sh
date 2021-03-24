#!/bin/bash

kubectl exec -tti $(kubectl get pods  -l "app=mattermost" -o jsonpath='{ .items[0].metadata.name }') -- sh