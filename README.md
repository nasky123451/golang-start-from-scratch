# Golang start from scratch  

If you want to use air, please make sure your go version >=1.23   

## 目录
  - [單元](#單元)
  - [指令](#指令)

### 單元

#### Websocket Unit

1. Go to folder

```   
cd websocket  
```   
2. Open Websocket Server (8080Port)  
```  
go run websocket_server.go   
```  
##### Or
``` 
go run websocket_server.go -monitor 
``` 
3. Open Websocket Client Test
Single Client Test 
``` 
go run websocket_client.go  
```   
##### Or Brute force test   
```   
go run websocket_clients.go   
```   

#### Tracing-Jaeger Unit

1. Go to folder

```   
cd tracing  
```   
2. Run Jaeger Server (16686Port)  

```   
docker run -d --name jaeger `
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
``` 

##### If you encounter

```   
docker: Error response from daemon: Conflict. The container name "/jaeger" is already in use by container "7e4c51a680ae4d5098fbc1b87d070229bceb459219feb7bc56eb2edce5f6d4d7". You have to remove (or rename) that container to be able to reuse that name.
See 'docker run --help'.
```   

##### Please run

```   
docker stop jaeger
docker rm jaeger
``` 

3. Go to browser

http://localhost:16686/   

4. Stop Jaeger Server (16686Port)  

```   
docker stop jaeger
``` 

#### Tracing-Zipkin Unit

1. Go to folder

```   
cd tracing  
```   
2. Run Zipkin Server (9412Port)  

```   
docker run -d --name zipkin -p 9412:9411 openzipkin/zipkin
``` 

3. Go to browser

http://localhost:9412/   

4. Stop Zipkin Server (9412Port)  

```   
docker stop zipkin
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
docker run -p 8080:8080 go-docker:latest  
# or   
docker run -p 8080:8080 --rm -v ${PWD}:/app -v /app/tmp --name go-docker-air go-docker

```   
