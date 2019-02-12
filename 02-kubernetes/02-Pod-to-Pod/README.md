## Pod to Pod communication

```
kubectl create -f employee.yaml

kubectl create -f stock-options.yaml

kubectl get all

kubectl describe service employee-service

```

## Access the employee service

http://MINIKUBE_IP:32200/employee


### Clean the setup

```
	kubectl delete -f ./
```