# in-toto-webhook

an in-toto admission controller for Kubernetes pods. This is still a work in
progress, so use at your own risk!

### Install

Generate webhook configuration files with a new TLS certificate and CA Bundle:

```bash
make certs
```

Deploy the admission controller and webhooks in the in-toto namespace (requires Kubernetes 1.10 or newer):

```bash
make deploy
``` 

Enable Kubesec validation by adding this label:

```bash
kubectl label namespaces default in-toto-validation=enabled
```

### Usage

Try to apply a privileged Deployment:

```bash
kubectl apply -f ./test/deployment.yaml

Error from server (InternalError): error when creating "./test/deployment.yaml": 
Internal error occurred: admission webhook "deployment.admission.kubesc.io" denied the request: 
deployment-test score is -30, deployment minimum accepted score is 0
```

Try to apply a privileged DaemonSet:

```bash
kubectl apply -f ./test/daemonset.yaml

Error from server (InternalError): error when creating "./test/daemonset.yaml": 
Internal error occurred: admission webhook "daemonset.admission.kubesc.io" denied the request: 
daemonset-test score is -30, daemonset minimum accepted score is 0
```

Try to apply a privileged StatefulSet:

```bash
kubectl apply -f ./test/statefulset.yaml

Error from server (InternalError): error when creating "./test/statefulset.yaml": 
Internal error occurred: admission webhook "statefulset.admission.kubesc.io" denied the request: 
statefulset-test score is -30, deployment minimum accepted score is 0
```

### Configuration

You can set the minimum Kubesec.io score in `./deploy/webhook/yaml`:

```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: in-toto-webhook
  labels:
    app: in-toto-webhook
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: in-toto-webhook
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8081"
    spec:
      containers:
        - name: in-toto-webhook
          image: stefanprodan/in-toto:0.1-dev
          imagePullPolicy: Always
          command:
            - ./in-toto
          args:
            - -tls-cert-file=/etc/webhook/certs/cert.pem
            - -tls-key-file=/etc/webhook/certs/key.pem
            - -min-score=0
          ports:
            - containerPort: 8080
            - containerPort: 8081
          volumeMounts:
            - name: webhook-certs
              mountPath: /etc/webhook/certs
              readOnly: true
      volumes:
        - name: webhook-certs
          secret:
            secretName: in-toto-webhook-certs
```

### serializig in-toto metadata

You can serialize in-toto metadata by doing a POST request to
https://yoururl:8000/links/namespace/linkname and a payload with the actual
link metadata. You can test it with make testserialization.

### Monitoring 

The admission controller exposes Prometheus RED metrics for each webhook a Grafana dashboard is available [here](https://grafana.com/dashboards/7088).

### Credits

This was *very heavily* based of [stefanprodan's kubesec webhook](https://github.com/stefanprodan/kubesec-webhook)

Kudos to [Xabier](https://github.com/slok) for the awesome [kubewebhook library](https://github.com/slok/kubewebhook).  
