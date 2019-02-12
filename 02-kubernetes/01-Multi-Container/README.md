## Multiple Containers in a Pod

```
kubectl create -f multi-container-pod.yaml

kubectl get pods

kubectl describe pods multi-container-pod

kubectl exec multi-container-pod -c 1st -- /bin/cat /usr/share/nginx/html/index.html

kubectl exec multi-container-pod -c 2nd -- /bin/cat /html/index.html

```


### Clean the setup

```
	kubectl delete -f multi-container-pod.yaml
```