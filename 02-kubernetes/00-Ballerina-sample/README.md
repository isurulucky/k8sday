## Deploy Ballerina Sample APP in K8s

```
kubectl create -f ballerina-sample.yaml

kubectl get all

kubectl get deployments

kubectl get pods

kubectl describe pod

kubectl describe deployment

kubectl describe replicaset

kubectl delete pod POD_ID

kubectl logs -f POD_ID

kubectl exec -it POD_ID bash

kubectl get services

kubectl describe services ballerina-service

```

### Access the service via NodePort

```
http://MINIKUBE_IP:32100/helloWorld/sayHello
```

### Deploy the ingress controller

```
minikube addons enable ingress

kubectl create -f ingress.yaml

kubectl get ing

```

### Update /etc/hosts file

```
IP_Address k8sday.com
```

### Access the service via Ingress

```
http://k8sday/helloWorld/sayHello
```

### K8s Labels

```
kubectl get pods -l event=k8sday

kubectl describe pod POD_ID

kubectl describe service ballerina-service

kubectl label pod POD_ID event-

```

### Readiness and Liveness Probes


### Clean the setup

```
kubectl delete -f ./
```

