# base-image:tag 
FROM golang:alpine
# Try to choose the smallest possible image size

# copies everything in current directory to /app in container. 
COPY . /app  
# Files listed in .dockerignore will be ignored. 

# To copy only specific files 
# COPY package.json /app

# A trailing forward slash is needed in destination (/app/) if multiple files are copied. 
# COPY package.json package-lock.json /app/

# Or use a pattern (will copy any matching file)
# COPY package*.json /app 

# sets workign directory to /app in container 
WORKDIR /app

# Destination paths are relative to WORKING DIRECTORY
# COPY README.md . 

# Use this syntax for file names with spaces (last element is destination)
# COPY ["hello world.txt", "."]

# Add is similar to COPY with additional features 
# ADD http://somedomain.com/file.json .  
# Like fetching files from URL or uncompressing zip. 
# ADD file.zip .
# Contents will be extracted into WORKDIR 


# The RUN commands executes any commands that are normally run in terminal 
# RUN npm install 
# RUN apt install python 

# Use ENV to se environment variables 
ENV API_URL=http://api.myapp.com/ 

# EXPOSE describes what port this container will be listening to. Exposes port for mapping later. 
EXPOSE 8080

# To set the current user (default is root). 
# USER app 
# All the following commands will be run by this user. Preceeding commands will still be run as default user (root)

# User app must be already created (can use RUN to create one before this)
# RUN addgroup app && adduser -S -G app app 


# You can run command in container by starting it like: 
# > docker run image-name npm start 
# OR 
# > docker run go-docker go run ./server/server.go 

# To avoid specifying run command everytime we run contianer, use CMD 
# With CMD we can supply a default command to be executed. 
CMD go run ./server/server.go
# Now we can start container by simply:
# > docker run go-docker 

# Only one CMD gets executed (last one by default)
# CMD npm start 

# RUN vs CMD > RUN is excuted when image is built (in intermediate container)
# > CMD is executed when starting a container (from the built image) 
# CMD > run-time instruction vs RUN > build-time instruction

# CMD has two forms 
# 1. Shell form 
# CMD npm start 
# Shell form is executed in a separate shell. 
# Basically: RUN /bin/sh npm start

# 2. Exec form (Prefer this)
# CMD ["npm", "start"]
# Exec from executes in current shell. No new shell is spin up for this. 

# CMD can always still be overidden when starting the container: 
# > docker run go-docker echo hello 
# "echo hello" will run instead of CMD (npm start)

# ENTRYPOINT is similar to CMD except it cannot be easily overridden
# ENTRYPOINT ["npm",  "start"]

# To override ENTRYPOINT --entrypoint option needs to be used. 
# > docker run go-docker --entrypoint echo hello 
# Avoids accidental overrides I guess. 