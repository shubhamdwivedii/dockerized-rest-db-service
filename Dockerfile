FROM golang:alpine 
COPY . /app  
WORKDIR /app
CMD go run ./server/server.go
EXPOSE 8080