# Golang start from scratch  

If you want to use air, please make sure your go version >=1.23   

If you want to use Prometheus, please download it from [Prometheus下載](https://prometheus.io/download/#:~:text=An%20open-source%20monitoring%20system%20with%20a) place
1. Unzip the downloaded folder
2. Copy prometheus.exe to %GOROOT%\bin\

You need to set up a Docker network so that your application and Zipkin and Jaeger containers can communicate with each other.
```  
docker network create my-network
```  

## 目录
  - [單元](#單元)
    - [Goroutine](#Goroutine Unit)
    - [Websocket](#Websocket Unit)
    - [Tracing](#Tracing Unit)
    - [Prometheus](#Prometheus Unit)
  - [指令](#指令)
    - [Git](#Git common commands)
    - [Docker](#Docker common commands)

### 單元

You can use -help to query tags
``` 
go run .\main.go -help
# or 
docker run --rm --name go-docker go-docker:latest -help  
``` 
#### Goroutine Unit

##### Goroutine

``` 
go run .\main.go -goroutine
# or 
docker run --rm --name go-docker go-docker:latest -goroutine
```

##### Goroutine Mutex

``` 
go run .\main.go -goroutineMutex
# or  
docker run --rm --name go-docker go-docker:latest -goroutineMutex
```

##### Goroutine Channel

``` 
go run .\main.go -goroutineChannel
# or  
docker run --rm --name go-docker go-docker:latest -goroutineChannel
```

#### Websocket Unit

##### Server (8080 Port)

```   
go run .\main.go -websocketServer
# or   
go run .\main.go -websocketServer -monitor
# or   
docker run --rm --name go-docker -p 8080:8080 go-docker:latest -websocketServer
# or   
docker run --rm --name go-docker -p 8080:8080 go-docker:latest -websocketServer -monitor
``` 

##### Client

```   
go run .\main.go -websocketClient
# or   
docker run --rm --name go-docker -p 8080:8080 go-docker:latest -websocketClient
``` 

##### Clients

```   
go run .\main.go -websocketClients
# or   
docker run --rm --name go-docker -p 8080:8080 go-docker:latest -websocketClients
``` 

#### Tracing Unit

##### Jaeger

1. Run Jaeger Server (16686Port)  

```   
docker run -d --rm --name jaeger `
  --network my-network `
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 `
  -p 5775:5775/udp `
  -p 6831:6831/udp `
  -p 6832:6832/udp `
  -p 5778:5778 `
  -p 16686:16686 `
  -p 14268:14268 `
  -p 14250:14250 `
  -p 9411:9411 `
  jaegertracing/all-in-one:1.32

docker run --rm --name go-docker --network my-network go-docker:latest -tracingJeager
``` 

2. Go to browser

http://localhost:16686/   

3. Stop Jaeger Server (16686Port)  

```   
docker stop jaeger
``` 

##### Zipkin

1. Run Zipkin Server (9412Port)  

```   
docker run -d --rm --name zipkin --network my-network -p 9412:9411 openzipkin/zipkin
docker run --rm --name go-docker --network my-network go-docker:latest -tracingZipkin
``` 

2. Go to browser

http://localhost:9412/   

3. Stop Zipkin Server (9412Port)  

```   
docker stop zipkin
``` 

#### Prometheus Unit

##### Prometheus

```   
go run .\main.go -prometheus
# or   
docker run --rm --name go-docker -p 8080:8080 -p 9090:9090 go-docker:latest -prometheus
``` 

### 指令

#### Git common commands
``` 
git add .   
git commit -m "Init"   
git push -u origin main   
``` 

#### Docker common commands
```   
docker build -t go-docker:latest .   
docker images 
docker run --rm --name go-docker go-docker:latest  
# or   
docker run --rm --name go-docker -v ${PWD}:/app -v /app/tmp --name go-docker-air go-docker

```   
#### Docker stop commands
```   
docker ps
docker stop go-docker
```   

#### If you encounter   

##### docker: Error response from daemon: Conflict. The container name "/xxxxxxxx" is already in use by container "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx". You have to remove (or rename) that container to be able to reuse that name.
##### See 'docker run --help'.

##### Please run

```   
docker stop xxxxxxxx
docker rm xxxxxxxx
``` 