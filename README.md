# Matrics Scrapper

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

 
```
