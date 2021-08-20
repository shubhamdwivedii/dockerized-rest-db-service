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

// DOCKER CONTAINERS ###########

// -d stands for "detached mode" ie: container will run in background.
// > docker run -d my-image

// Use --name to give container a name when starting
// > docker run -d --name my-container my-image

// To see logs on container
// > docker logs container-id

// Only first few letters of continer-ids (or even image-ids) are needed (as long as they are unique)
// Eg for continaer id 65529e445cbc (if not other containers start with 655)
// > docker logs 655

// If contanier is continuously logging use -f or --follow to follow logs in real time
// > docker logs -f 655

// Use -n or --tail to specify last n lines of logs to print
// > docker logs -n 5 655

// Use -t to also print timestams of logs
// > docker logs -n 8 -t 655

// To publish port (map container port to host port)
// > docker run -d -p 8080:3000 --name my-container my-image
// 8080 is host port : 3000 is container port (server is listening to 3000 in container for example)

// To run a command on a running contianer (later on)
// > docker exec container-one ls
// Note this will run in WORKDIR set in Dockerfile

// NOTE: below "my-container" is container name (container id can be used also)

// We can even start a shell session
// > docker exec -it my-container sh

// "run" runs (starts) a container & "exec" executes a command in running container

// To execute as a root user (or any other) use -u
// > docker exec -it -u root my-container sh

// To stop container
// > docker stop my-container

// To restart a stopped container
// > docker start my-continaer

// Note: use "run" to start new containers and "start" to restart stopped containers

// To remove contianer
// > docker container rm my-container
// OR shortcut
// > docker rm my-container

// You cannot remove running containers. They must be either stopped or use -f force option
// > docker rm -f my-container

// Removed container won't show in stopped container list (as they are removed duh!). TO verify
// > docker ps -a | grep my-container

// To remove all stopped containers
// > docker container prune

// CONTAINER FILE SYSTEM ############

// Each continaer has its own separate file system

// To verify start two containers with same image:
// > docker run -d --name c-one my-image
// > docker run -d --name c-two my-image

// Add a file in one container
// > docker exec -it c-one sh
// $ echo data > data.txt
// $ exit

// > docker exec -it c-two sh
// $ ls | grep data

// data.txt is not there in c-two.

// Each container has its own file system. If we delete c-one, we'll lose all its files including data.txt
// Since we might lost important data when deleting continers,
// We should NEVER store our data in a container's file system.

// Persisting data with VOLUMES ##########

// A Volume is a storage outside of containers. It could be a directory on host, or a bucket etc on cloud.

// Create a new Volume
// > docker volume create app-data
// > docker volume inspect app-data
/* {
	"CreatedAt": "2021-03016T21:16:41Z",
	"Driver": "local", // volume is directory on host.
	"Labesl": {},
	"Mountpoint": "C:\Shubham\DockerTut",
	"Name": "app-data",
	"Options": {},
	"Scope": "local",
    }
*/

// Start a container with this volume to persist data (using -v option)
// > docker run -d -p 4000:3000 -v app-data:/app/data my-image

// The volume is mapped to the /app/data directory inside the container.

// Volume don't need to be created before-hand
// > docker run -d -v new-volume:/app/data my-image
// This will work too.

// NOTE: if /data is not already created in /app in container,
// Docker will automatically create /data folder but as a "root" user.
// the app user we delcared in Dockerfile might not have access to /data folder this way.

// ALWAYS make sure any directory mapped to a Volume is already created with correct user in Dockerfile.

// If we delete the container now, any files created in /app/data will still exist in the Volume.

// Volumes can be SHARED among MULTIPLE Containers.

// COPYING Files b/w Containers ########

// Will copy /app/log.txt to current directory in Host (.)
// > docker cp source-container:/app/log.txt .

// Note: /app/log.txt is absolute path of file in container

// Copy secret.txt from Host into container's /app folder
// > docker cp secret.txt target-container:/app

// Sharing SOURCE with a Container #######

// Until now whenever source code changes, we create new docker image.

// In PRODUCTION > Always build a new image to publish changes in source.

// In DEVELOPMENT > We can create a mapping/binding b/w a directory on host & a directory inside a container.

// To create a mapping we can use Volume syntax but instead of a named volume we provide a directory in Host.
// > docker run -d -p 4000:3000 -v $(pwd):/app

// We can still also use a named Volume along with host directory mapping.
// > docker run -d -p 4000:3000 -v $(pwd):/app -v my-volume:/app/data my-image

// Any changes in the Host directory (pwd) will be immediately reflected in the Container.
// Example: Hot-Reloading in React can be achieved via this.

// To only see ids of images add -q
// > docker images ls -q

// This can be used as shortcut to delete all images at once
// > docker image rm $(docker image ls -q)

// You'll get error if some of these images are used in running or even stopped containers.

// Remove all containers first (use -f for force remove running containers)
// > docker container rm -f $(docker container ls -aq)

// To kill all containers
// > docker container kill $(docker ps -aq)
