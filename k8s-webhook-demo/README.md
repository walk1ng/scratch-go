# Kubernetes Webhook Demo

## validate 
> The replicas of deployment should ge 3.

## mutate
> The nginx container will be injected.

## deployment 
```bash
# sidecar container
kubectl create configmap webhook-sidecar --from-file manifests/sidecar.yaml

# service
kubectl apply -f config/samples/service.yaml

# csr and webhook tls
scripts/webhook-create-signed-cert.sh --service admission-webhook --namespace dev --secret webhook-certs

# webhook deployment
kubectl apply -f config/samples/deployment.yaml

# replace caBundle in webhook configuration
# get ca by script webhook-patch-ca-bundle.sh

# mutating webhook configuration
kubectl apply -f config/samples/mutating.yaml

# validating webhook configuration
kubectl apply -f config/samples/validating.yaml

# test app 
kubectl apply -f config/samples/testapp.yaml

```