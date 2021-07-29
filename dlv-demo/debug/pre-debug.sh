#!/bin/bash

go mod vendor
docker build -t my-golang-app-image -f debug/Dockerfile .
kubectl apply -f debug/debug.yaml

while :
do
  status=$(kubectl get pods | grep my-golang-app | awk '{print $3}')
  echo "pod status:" $status
  if [ "$status"  = "Running" ]; then
   echo "pod status checked pass"
   break
  fi
  sleep 5
done

while :
do
  line_num=$(kubectl get pods | grep my-golang-app | awk '{print $1}' | xargs kubectl logs | wc -l)
  if (( $line_num > 1 )); then
    echo "dlv server started"
    break
  fi
  sleep 5
done