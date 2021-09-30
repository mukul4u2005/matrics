# Metrics Scrapper

This is Prometheus Integration with rest service exposed using golang application. Golang application is deployed on kubernetes cluster and exposed as rest endpoint url.

There are two matrics exposed:

1) Gauge (To count number of pods in cluster)
2) Counter (Couting the total request hit count)

Pre-requisite:

## Create RBAC rules to allow go application to get pod list on cluster level.


```
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: pod-reader
rules:
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - pods

```

```
kind: ServiceAccount
apiVersion: v1
metadata:
  name: pod-service-account
  namespace: <namespace-name>

```

```
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: pod-reader-global
subjects:
  - kind: ServiceAccount
    name: pod-service-account
    namespace: mukul-test
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: pod-reader

```

# Build Metrics Go application, create a docker image and deploy on kubernetes

Clone code first:

```
git clone https://github.com/mukul4u2005/matrics.git
 
```

Build code:

```
GOOS=linux go build -o ./matrix-app .
```

Build docker image:

```
docker build -t matrix-app .
```
Tag image

```
docker tag matrix-app mpdocker2017/matrix-app:latest
```
Push image to Dockerhub 

```
docker push mpdocker2017/matrix-app:latest   
```

# Deploy Application on Kubernetes cluster and expose

Create Secret for cluster api call, this secret is referenced in metrics deployment to configure as env. variables:

```
 kubectl create secret generic cluster-secret --from-literal=TOKEN=<token> --from-literal=API_URL=<api-url>

```

Create deployment using yaml available in config folder:

For metrics service (Please change the image url):

```
kubectl apply -f metrics-deployment.yml

```
Create Prometheus config map for scrapper configuration (For prometheus deployment):

```
kubectl create configmap prome-config --from-file=<path>/prometheus.yml

```

Create Prometheus deployment:

```
kubectl apply -f prom-deployment.yml
```

# Access Prometheus using port-forward

kubectl port-forward pod-name 9090:9090


