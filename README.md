# Motivation

This project is a large collection of developers' test applications for various functions and suites of golang.

## 目录
  - [開發者的配置](#開發者的配置)
  - [使用方法](#使用方法)
  - [單元](#單元)
    - [Goroutine](#Goroutine)
    - [Websocket](#Websocket)
    - [Tracing](#Tracing)
    - [Prometheus](#Prometheus)
    - [Redis](#Redis)
  - [指令](#指令)
    - [Git](#Git)
    - [Docker](#Docker)
  - [常見問題](#常見問題)

## 開發者的配置

1. go - go1.23.1 windows/amd64
2. Docker - Docker version 27.2.0, build 3ab4256
3. git - git version 2.46.2.windows.1
4. air - air v1.60.0, built with Go go1.23.1
5. prometheus - prometheus, version 3.0.0-beta.0
6. PostgreSQL - PostgreSQL 17.0 on x86_64-windows, compiled by msvc-19.41.34120, 64-bit
7. Redis - Redis server v=5.0.14.1

## 使用方法

Welcome to the project! Follow these steps to get up and running quickly.

### Displaying All Test Functions
To view all available test functions, use the -help flag in the root directory:
``` 
go run .\main.go -help
``` 

### Running with Docker
This project supports Docker, and you can choose to run it using either Go commands or Docker commands:
```  
docker build -t go-docker:latest .
docker run --rm --name go-docker go-docker:latest -help
```  

### Setting Up Prometheus
If you'd like to integrate Prometheus for monitoring, follow these steps:

1. Download Prometheus from the [Prometheus Download page](https://prometheus.io/download/#:~:text=An%20open-source%20monitoring%20system%20with%20a).
2. Unzip the downloaded folder.
3. Copy prometheus.exe to the %GOROOT%\bin\ directory for easy access from the command line.

### Setting Up Redis
To use Redis in your project, download the Windows-compatible version from the [Redis for Windows Download page](https://github.com/tporadowski/redis/releases).

### Installing Air for Hot Reloading
For more efficient development and debugging, it's recommended to install the Air hot-reloading tool. Ensure your Go version is 1.23 or higher:
```  
go install github.com/air-verse/air@latest
air -v
```  

Configuring Docker Network
To enable communication between your application and Zipkin or Jaeger containers, create a Docker network:
```  
docker network create my-network
```  

## 單元

### Goroutine

  - Goroutine Base: Product inventory management
    - Function: Multiple consumers try to purchase goods and manage inventory through atomic operations.
    - Key point: Use atomic to safely modify the inventory and ensure data consistency.
  - Goroutine Mutex: Bank account operations
    - Function: Simulate a bank account and randomly perform deposit and withdrawal operations.
    - Key takeaway: Use sync.Mutex to ensure safe access to shared balances and avoid race conditions.
  - Goroutine Channel: Task producers and consumers
    - Function: Use Goroutine to generate random tasks and pass them to consumers for processing through channels.
    - Key Point: Demonstrates the producer-consumer pattern and how to use stop channel to end production.   

These examples demonstrate concurrent programming techniques in Go and are suitable for different application scenarios.

#### Goroutine Base

``` 
go run .\main.go -goroutine

# Run using docker
docker run --rm --name go-docker go-docker:latest -goroutine
```

#### Goroutine Mutex

``` 
go run .\main.go -goroutineMutex

# Run using docker 
docker run --rm --name go-docker go-docker:latest -goroutineMutex
```

#### Goroutine Channel

``` 
go run .\main.go -goroutineChannel

# Run using docker 
docker run --rm --name go-docker go-docker:latest -goroutineChannel
```

### Websocket

  - Server: WebSocket Hub server
    - Client: Represents the connected WebSocket client, including connection, sending channel and client ID.
    - Hub: Manages all connected clients, handling registration, deregistration and broadcast messages.
    - Function:
      - Use mutex locks to protect shared resources.
      - Handles client read and write operations and periodically sends Ping messages to maintain the connection.
      - Monitor and log system resource usage.
  - Client: WebSocket client
    - Function:
      - Connect to the WebSocket server and receive messages.
      - Use goroutine to process messages received from the server.
      - Send messages regularly and close the connection when completed.
  - Clients: WebSocket test client
    - TestClient: Test client that can connect to the server and send random messages.
    - Function:
      - Set heartbeat messages to stay connected.
      - Generate messages randomly and send them, closing the connection safely when done.
      - Supports concurrent testing of multiple clients.   

These examples show how to implement basic functionality of a WebSocket server and client using the Go language and the Gorilla WebSocket suite.

#### Server (8080 Port)

```   
go run .\main.go -websocketServer

# using monitor  
go run .\main.go -websocketServerMonitor

# Run using docker  
docker run --rm --name go-docker go-docker:latest -websocketServer

# Run using docker and using monitor 
docker run --rm --name go-docker go-docker:latest -websocketServerMonitor
``` 

#### Client

```   
go run .\main.go -websocketClient

# Run using docker  
docker run --rm --name go-docker go-docker:latest -websocketClient
``` 

#### Clients

```   
go run .\main.go -websocketClients

# Run using docker  
docker run --rm --name go-docker go-docker:latest -websocketClients
``` 

### Tracing

  - Jaeger: Jaeger tracking
    - Purpose: Use Jaeger for distributed tracing to help monitor and debug microservice architecture.
    - Key features:
      - Initialize the Jaeger tracer: Set up the exporter and tracer provider.
      - Create and end trace spans: Use tracer.Start and span.End() to trace operations.
      - Data export: Make sure the exporter has time to send the data to Jaeger.
  - Zipkin: Zipkin tracking
    - Purpose: Use Zipkin for distributed tracing, similar to Jaeger.
    - Key features:
      - Initialize the Zipkin tracker: Set up the exporter and tracker provider.
      - Creating and ending trace spans: Also use tracer.Start and span.End().
      - Data export: Make sure the exporter has time to send data to Zipkin.   

These two examples show how to use OpenTelemetry to integrate Jaeger and Zipkin for distributed tracing to help analyze the performance and request flow of microservices.

#### Jaeger

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

docker run --rm --name go-docker --network my-network -e "URL=jaeger" go-docker:latest -tracingJeager
``` 

2. Go to browser

http://localhost:16686/   

3. Stop Jaeger Server (16686Port)  

```   
docker stop jaeger
``` 

#### Zipkin

1. Run Zipkin Server (9412Port)  

``` 
docker run -d --rm --name zipkin --network my-network -p 9412:9411 openzipkin/zipkin  
docker run --rm --name go-docker --network my-network -e "URL=zipkin" go-docker:latest -tracingZipkin
``` 

2. Go to browser

http://localhost:9412/   

3. Stop Zipkin Server (9412Port)  

```   
docker stop zipkin
``` 

### Prometheus

  - Prometheus Base: Basic Prometheus Application
    - Purpose: Sets up a simple HTTP server that integrates with Prometheus for basic metrics.
    - Key Features:
      - Responds with "Hello, World!" at the root path ("/").
      - Exposes a /metrics endpoint for Prometheus to scrape metrics.
      - Tracks request count and duration for monitoring.
  - Prometheus API Application: Prometheus API Application with Database
    - Purpose: Expands the first example to include database interactions and multiple API routes.
    - Key Features:
      - Listens on port 8080 with multiple routes:
        - /api/v1/resource: Fetches and returns resources from a PostgreSQL database as JSON.
        - /api/v1/login: Simulates a login response.
        - /health: Returns a health check status.
        - /metrics: Serves Prometheus metrics.
      - Supports graceful shutdown, allowing cleanup before terminating.
      - Updates metrics for request counts, durations, and latencies.   

Both examples demonstrate how to integrate Prometheus into a Go application, with the first being a simple server and the second providing a more complex API with database functionality.

#### Prometheus Base

```   
go run .\main.go -prometheus

# Run using docker  
docker run --rm --name go-docker go-docker:latest -prometheus
``` 

#### Prometheus API Application

1. Run Postgres Server (5432Port)  

```   
go run .\main.go -prometheusApiApplication

# Run using docker  
docker run -d --rm --name postgres-container --network my-network -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=henry -e POSTGRES_DB=test postgres:latest
docker run --rm --name go-docker --network my-network -e "DATABASE_URL=postgres-container" go-docker:latest -prometheusApiApplication
``` 

2. Go to browser

http://localhost:9090/   

3. Stop Postgres Server (5432Port)  

```   
docker stop postgres-container
``` 

### Redis

Under development

#### Redis Base

1. Run Postgres Server (5432Port & 6379port)  

```   
go run .\main.go -redisbase

# Run using docker  
docker run -d --rm --name redis-container --network my-network redis:latest
docker run -d --rm --name postgres-container --network my-network -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=henry -e POSTGRES_DB=test postgres:latest
docker run --rm --name go-docker --network my-network -e "DATABASE_URL=postgres-container" -e "REDIS_URL=redis-container" go-docker:latest -redisbase
``` 

2. Stop Postgres Server (5432Port & 6379port)  

```   
docker stop redis-container
docker stop postgres-container
``` 

## 指令

### Git

Here is a record of commonly used commands in Git

#### Git common commands
``` 
git add .   
git commit -m "Init"   
git push -u origin main   
``` 

### Docker

Here is a record of commonly used commands in Docker

#### Docker common commands
```   
docker build -t go-docker:latest .   
docker images 
docker run --rm --name go-docker go-docker:latest  

# Run using docker  
docker run --rm --name go-docker -v ${PWD}:/app -v /app/tmp --name go-docker-air go-docker

```   
### Docker stop commands
```   
docker ps
docker stop <NAMES>
```   

## 常見問題

#### docker: Error response from daemon: Conflict. The container name "/xxxxxxxx" is already in use by container "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx". You have to remove (or rename) that container to be able to reuse that name.
#### See 'docker run --help'.

  - Please run
```   
docker stop xxxxxxxx
docker rm xxxxxxxx
``` 

### "command not found: air" or "No such file or directory"

  - You can also directly copy the exe file to %GOROOT%\bin\, or use the following command
```   
export GOPATH=$HOME/xxxxx
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export PATH=$PATH:$(go env GOPATH)/bin <---- Confirm this line in you profile!!!
```   