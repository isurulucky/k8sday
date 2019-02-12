## Performing Rolling Update

```
kubectl create -f nginx.yaml

kubectl rolling-update my-nginx --image=nginx:1.9.1

kubectl describe pod

```

### Perform rollback 

```
kubectl rolling-update my-nginx --rollback
```

### Clean the setup

```
kubectl delete -f ./
```
