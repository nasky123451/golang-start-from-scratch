# Golang start from scratch  

If you want to use air, please make sure your go version >=1.23   

# Websocket Unit

1. Go to folder

```   
cd websocket  
```   
2. Open Websocket Server (8080Port)  
```  
go run websocket_server.go   
```  
Or
``` 
go run websocket_server.go -monitor 
``` 
3. Open Websocket Client Test
Single Client Test 
``` 
go run websocket_client.go  
```   
Or Brute force test   
```   
go run websocket_clients.go   
```   

# Git common commands
``` 
git add .   
git commit -m "Init"   
git push -u origin main   
``` 

# Docker common commands
```   
docker build -t go-docker:latest .   
docker images 
docker run -it go-docker:latest air init  
docker run -p 8080:8080 go-docker:latest   
```   
