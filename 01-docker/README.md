# Agenda
* Working with Containers
* Working with Images
* Docker best practices
* Additional content

## Working with Containers
* Pull docker image from docker hub(public docker registry)
```
docker pull manjula/k8sday.demo.service:1.0.0
```
* View pulled docker images in your docker host(local machine)
```
docker images
```
* Run a container in foreground using the docker image pulled to your docker host
```
docker run manjula/k8sday.demo.service:1.0.0
```
* View running docker containers in your docker host
```
docker ps
```
* Run another container in background and verify with ```docker ps```
```
docker run -d manjula/k8sday.demo.service:1.0.0
```
* View processes on docker container
```
docker exec -it containerID ps 
```
* View processes on docker host machine
```
ps aux | grep hello_service 
```
* Kill processes on docker host machine and notice that container processes are killed too. ```docker ps``` 
```
sudo kill -9 processID of /bin/sh -c ballerina run  hello_service.balx
```
* Run a command inside the running container. 
```
docker exec -it 5b pwd
docker exec -it 5b /bin/bash 
```
* Invoke the ballerina service within the container. 
```
docker exec -it 5b /bin/bash 
wget -O- http://127.0.0.1:9090/helloWorld/sayHello
```
* Invoke the ballerina service from the docker host itself. You need to map the docker host port to the container port
```
docker run -d -p 8000:9090 hello_service:latest
docker run -d -p 8000:9090 hello_service:latest (notice the port binding error)
docker run -d -p 9000:9090 hello_service:latest
docker run -d -P hello_service:latest (check with docker ps and notice the host port is dynamically picked and mapped to the exported container port)
```
Invoke above services using the docker host.
```
curl http://127.0.0.1:8000/helloWorld/sayHello
curl http://127.0.0.1:9000/helloWorld/sayHellom
curl http://127.0.0.1:32768/helloWorld/sayHello
```
* View other details of the running containers.
```
docker logs 5b (shows logs of the container)
docker inspect 5b (shows the complete details such as docker image, envs, commands, mounts, network settings etc)
docker top 5b (shows running processes in container)
docker stat 5b (shows containers' resource usage statistics)
docker diff 5b (shows changed files in the container's file system)
```
* Run container with cpu, memory resource limits
```
nproc -all (get number of cpus in docker host)
docker run -it --cpus=1 --memory=100m manjula/k8sday.demo.service:1.0.0
```
* Mount files from docker host to the container
```
docker run -v /home/manjula/Desktop/demo/:/tmp/ -it  manjula/k8sday.demo.service:1.0.0
```
* Clean up stopped containers
```
docker rm -f c751e33514b0
docker rm -f $(docker ps -aq ) (remove all stopped containers)
docker run -d --rm manjula/k8sday.demo.service:1.0.0 (use --rm to destory the container when container is stopped/killed)
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
docker history manjula/k8sday.demo.service:1.0.0
```
* Remove docker images
```
docker rmi -f manjula/k8sday.demo.service:1.0.0
```

## DockerBest practices
* Container should be immutable. No patching of running container, rather start with a new image with patches and run a container of it.
* Do not run software as root. Create a less privileged user.
* No need to install ssh iside the container, use exec -it image /bin/bash instead.
* Run containers with resource limts avoid taking complete resouces of the host machine.
* Container should have only one process 
* Reduce the services installed on the docker image so that the size of the image is reduced and the security of the image is increased due to less number of softwares installed.
* Dont store senstive data inside the image, instead pass them through volume mounts or as envs.
* Follow best practices when writing docker files - https://docs.docker.com/develop/develop-images/dockerfile_best-practices/

	
## Additional content
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
