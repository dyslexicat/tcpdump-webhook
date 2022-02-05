# tcpdump-webhook

A simple demonstration of Kubernetes Mutating Webhooks. Injects a tcpdump sidecar to your pods with the `tcpdump-sidecar` label.

Rewrite of https://github.com/bilalunalnet/k8s-tcpdump-webhook in **Golang**

DISCLAIMER: *THIS IS NOT MEANT FOR PRODUCTION*

## Build and Deploy

In order to get your sidecar containers running you need to build the following image too: https://github.com/dyslexicat/tcpdump-alpine

Build docker image;

`docker build -t dyslexicat/tcpdump-webhook .`

Generate ca in /tmp
`cfssl gencert -initca ./tls/ca-csr.json | cfssljson -bare /tmp/ca`


Generate private key and certificate for SSL connection.

```
cfssl gencert \
  -ca=/tmp/ca.pem \
  -ca-key=/tmp/ca-key.pem \
  -config=./tls/ca-config.json \
  -hostname="tcpdump-webhook,tcpdump-webhook.webhook-demo.svc.cluster.local,tcpdump-webhook.webhook-demo.svc,localhost,127.0.0.1" \
  -profile=default \
  ./tls/ca-csr.json | cfssljson -bare /tmp/tcpdump-webhook
```

Move your SSL key and certificate to the `ssl` directory:
`mv /tmp/tcpdump-webhook.pem ./ssl/tcpdump.pem`
`mv /tmp/tcpdump-webhook-key.pem ./ssl/tcpdump.key`

Update ConfigMap data in the `manifest/webhook-deployment.yaml` file with your key and certificate.

Update `caBundle` value in the `manifest/webhook-configuration.yaml` file with your base64 encoded CA certificate.

`cat ca.pem | base64`

```
kubectl create ns webhook-demo
kubectl apply -f manifest/webhook-deployment.yaml
kubectl apply -f manifest/webhook-configuration.yaml
```

## Test

There is a Pod manifest file in the `manifest` directory to be used for testing purposes. 
