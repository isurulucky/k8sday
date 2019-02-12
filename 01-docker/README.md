# Docker

## Building Docker Image
```docker
docker build -t hello_service:latest ./docker
```

## Running Docker contianer
```docker
docker run -d -p 9090:9090 hello_service:latest
```

## Access the container
```bash
curl http://127.0.0.1:9090/helloWorld/sayHello
```