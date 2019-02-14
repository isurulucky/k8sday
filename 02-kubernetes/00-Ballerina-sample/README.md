## Deploy Ballerina Sample APP in K8s

```
docker save pubudu/hello_service:latest > image.tar

minikube docker-env

docker load < image.tar

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
kubectl get service ballerina-service -o json

kubectl get service ballerina-service -o jsonpath='{.spec.ports[?(@.name=="web")].nodePort}'

http://MINIKUBE_IP:32100/helloWorld/sayHello
```

* Note: MINIKUBE_IP is localhost for Mac users

### Deploy the ingress controller

#### For Docker For Mac

```
kubectl apply -f dockerForMac/nginx-ingress/namespaces/nginx-ingress.yaml -Rf dockerForMac/nginx-ingress
```

#### For Minikube
```
minikube addons enable ingress
```

```
## For Minikube only
kubectl get pods -n kube-system

## For Mac only
kubectl get pods -n nginx-ingress

kubectl create -f ingress.yaml

kubectl get ing

```

### Update /etc/hosts file

```
IP_Address k8sday.com
```

* Note: IP_Address is localhost for Mac users

### Access the service via Ingress

```
http://k8sday.com/helloWorld/sayHello
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
kubectl delete -f ./ -R

kubectl delete pod POD_ID
```

