# Agenda
* Working with Containers
* Working with Images
* Additional content

## Working with Containers
* Run Hello container
```
docker run hello-world
```
* Run ballerina hello service container
```
cat hello_service.bal 
docker run pubudu/hello_service:latest
```
* View running docker containers in your docker host
```
docker ps
```

* Run another container in background/detached mode and verify with ```docker ps```
```
docker run -d pubudu/hello_service:latest 
```
* View processes on docker container
```
docker exec -it <ContainerID> ps 
```
* View processes on docker host machine
```
ps aux | grep hello_service 
```
* Kill processes on docker host machine and notice that container processes are killed too. ```docker ps``` 
```
sudo kill -9 <ProcessID> of /bin/sh -c ballerina run  hello_service.balx
```
* Run a command inside the running container. 
```
docker exec -it <ContainerID> pwd
docker exec -it <ContainerID> /bin/bash 
```
* Invoke the ballerina service within the container. 
```
docker exec -it <ContainerID> /bin/bash 
wget -O- http://127.0.0.1:9090/helloWorld/sayHello
```
* Invoke the ballerina service from the docker host itself. You need to map the docker host port to the container port. This will use the default bridge network
```
docker run -d -p 8000:9090 pubudu/hello_service:latest 
docker run -d -p 8000:9090 pubudu/hello_service:latest (notice the port binding error)
docker run -d -p 9000:9090 pubudu/hello_service:latest 
docker run -d -P pubudu/hello_service:latest (check with docker ps and notice the host port is dynamically picked and mapped to the exported container port)
```
Invoke above services using the docker host.
```
curl http://127.0.0.1:8000/helloWorld/sayHello
curl http://127.0.0.1:9000/helloWorld/sayHello
curl http://127.0.0.1:<HostPort>/helloWorld/sayHello
```
* View other details of the running containers.
```
docker logs <ContainerID> (shows logs of the container)
docker inspect <ContainerID> (shows the complete details such as docker image, envs, commands, mounts, network settings etc)
docker top <ContainerID> (shows running processes in container)
docker stat <ContainerID> (shows containers' resource usage statistics)
docker diff <ContainerID> (shows changed files in the container's file system)
```
* Run container with cpu, memory resource limits
```
docker run -it --cpus=1 --memory=100m pubudu/hello_service:latest 
```
* Mount files from docker host to the container
```
docker run -v /home/manjula/Desktop/demo/:/tmp/ -it pubudu/hello_service:latest 
```
* Clean up stopped containers
```
docker rm -f <ContainerID> 
docker rm -f $(docker ps -aq ) (remove all stopped containers)
docker run -d --rm pubudu/hello_service:latest (use --rm to destory the container when container is exited)
```

## Working with Images
* Build a docker image
```
docker build -t hello_service:latest ./docker
docker images
```
* Tag a docker image
```
docker tag 7b manjula/k8sday.demo.service:2.0.0
docker images ( notice that multiple tags are created for the same docker imageID)
```
* Push docker image to docker hub 
prerequisite: 
1. have a account in https://hub.docker.com/ and create a docker repository 'manjula/k8sday.demo.service'
2. login to docker registry from docker cli (```docker login ``` caution: your credetials are stored in /home/user/.docker/config file)
```
docker push manjula/k8sday.demo.service:2.0.0
```
* View different images that builds up the docker image
```
docker history pubudu/hello_service:latest 
```
* Remove docker images
```
docker rmi -f pubudu/hello_service:latest 
```
## Additional content
* To run docker commands without sudo, https://docs.docker.com/install/linux/linux-postinstall/

* Docker status
  ```
  service docker status
  ```
* Docker info such as version, container, image stats, docker root directory
  ```
  docker info
  ```
* Enable debug logs for Docker Server(Daemon)
  ```
    /etc/docker/daemon.json
		{
		 "experimental": true,
		 "debug": true
		}
	service docker reload
  ```
* View Docker server logs
  ```
  journalctl --follow -u docker.service
  ```
* Docker cheat sheet - https://github.com/wsargent/docker-cheat-sheet
