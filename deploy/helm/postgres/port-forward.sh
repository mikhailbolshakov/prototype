#!/bin/bash

# connect from outside
kubectl port-forward --namespace default svc/pg-postgresql 5432:5432