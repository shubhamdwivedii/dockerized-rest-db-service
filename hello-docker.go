package main

import "fmt"

func main() {
	fmt.Println("Hello Docker!!")
}

// To build image
// docker build -t go-docker .

// To run image
// docker run go-docker

// To pull and run image directly (from dockerhub)
// docker pull golang/go-image
// docker images
// docker run golang/go-image

// TO expose and map port
//  docker run -d -p 8080:8080 --name my-container go-docker
// my-container is name of container go-docker is image name
