# Hello CRD



#### Building controller

    go build -o hello-controller -x ./cmd/controller

#### Build and push to your own repository

1. Open the ./build.sh and change the `DOCKER_REPO` variable with your docker ID
2. Run ./build.sh

#### Running the controller in your machine

    # Here the cluster connection infomation is loaded from kubeconfig file
    ./hello-controller -logtostderr=true --kubeconfig=/home/<username>/.kube/config

#### Running the controller in the cluster

    kubectl apply -f ./artifacts/crd.yaml
    # If you change the docker repo, please change the controller.yaml to pull from your repo
    kubectl apply -f ./artifacts/controller.yaml


#### Running the sample

    kubectl apply -f ./sample/my-hello.yaml
    
#### Testing the sample

    # Create a pod to debug so that we can curl
    kubectl run debug-tools --image=mirage20/k8s-debug-tools --restart=Never
    
    # Exec into the pod
    kubectl exec -it debug-tools /bin/bash
    
    # Inside the pod run the following command
    curl my-hello-service
    
    # Yous should see an output similar to following
    Hello, World from my-hello-deployment-848dbfbd67-fghqg



If you have a cluster with Istio installed, please disable sidecar injection

    kubectl label namespace default istio-injection=disabled --overwrite


#### Cleaning the cluster

    kubectl delete -f ./sample/my-hello.yaml
    kubectl delete -f ./artifacts/controller.yaml
    kubectl delete -f ./artifacts/crd.yaml
    kubectl delete pod debug-tools
