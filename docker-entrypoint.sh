# !/bin/sh 

echo "Waiting for DB to start..."
ls 
./wait-for-it database:3306 
# wait-for-it is a custom script for docker. It stalls until DB is intialized at PORT. 
# wait-for-it is available at https://github.com/vishnubob/wait-for-it

echo "Migrating the database..."
# migrage db commands (in this example: will run in the Backend container (where this sh file is called))

echo "Starting the server..."
go run ./server/server.go

# REMEMBER: this will override the Dockerfile's CMD. 
# Thus server needs to be started here. 