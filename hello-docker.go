package main

import "fmt"

func main() {
	fmt.Println("Hello Docker!!")
}

// To build image
// > docker build -t go-docker .

// To build image and add a tag
// > docker build -t image-name:image-tag .
// > docker build -t go-docker:1.2 .

// To add tag later
// > docker image tag image-name:latest

// multiple tags can be added to same image
// > docker image tag image-name:first-tag image-name:second-tag

// To switch tag to different version of image
// > docker image tag image-id image-name:latest

// NOTE image-name is Repository.
// A Repository can have multiple images with differenct image-ids
// An image with image-id can have multiple tags. (like latest, v1.3)

// To run image
// > docker run go-docker

// To run image with a command (command is run when container is launched)
// > docker run node-docker npm start

// To pull and run image directly (from dockerhub)
// > docker pull golang/go-image
// > docker images
// > docker run golang/go-image

// TO expose and map port
// > docker run -d -p 8080:8080 --name my-container go-docker
// my-container is name of container go-docker is image name

// To see all images
// > docker images
// > docker image ls

// To see all running containers
// > docker ps

// To see all containers
// > docker ps -a

// To run in interactive mode (opens shell)
// > docker run -it image-name
// > docker run -it image-name bash
// > docker run -it image-name sh

// To see image history and all its layers
// > docker history image-name
// Layers are cached and reused when images are rebuilt

// NOTE: use COPY . /app after running COPY package.json and RUN npm start
// If image is rebuilt after changes in src, cache is not used for any commands after COPY . /app
// This way node_modules will still be reused if there are no changes in package.json
// This will speed up the builds

// To delete danglings images (unused image layers)
// > docker image prune

// To remove a specific image(s)
// > docker image rm image-name
// > docker image rm first-image second-image third-image

// > docker image remove image-name:img-tag
// If image has multiple tags, only the tag gets remove (otherwise the whole image)

// To delete stopped containers
// > docker container prune

// To stop a container
// > docker container kill container-tag
// > docker container kill container-name

// To push image to dockerhub
// > docker login
// > docker push username/image-name:tag
// > docker push username/repository:tag

// image-name is same as repository.

// Add username to image as:
// > docker image tag image-id username/image-name:tag-name

// By default latest tag is used automatically.

// To save image as compressed file
// > docker image save -o image-name.tar image-name:tag
// > docker image save -o my-image.zip image-name:latest

// To load a saved image
// > docker image load -i image-name.tar
// > docker image load -i my-image.zip
