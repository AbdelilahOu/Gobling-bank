# Go Backend Master Class

## Getting Started

Welcome to the Go Backend Master Class project. This guide will help you set up the development environment and get started with the project.

### Setting up a Database with Docker

We use PostgreSQL as the database for this project. You can set up a PostgreSQL database with Docker using the following steps:

1. **Pull the PostgreSQL Docker image**:

   ```bash
   docker pull postgres:15
   ```

2. **Start a PostgreSQL Instance**:

   ```bash
   docker run --name postgres-database-dev -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:15
   ```

3. **Check if your container is running**:

   Run the following command and look for your container name in the list:

   ```bash
   dcoker ps
   ```

### Generating Queries Using SQLC

We use SQLC for generating Go code from SQL queries, To get started with SQLC, and avoid any errors we will use Docker.

first this is how you execute any command in docker:

    ```bash
    docker exec -it <container_name_or_container_id> <command>  [args]
    ```

1. **Pull the SQLC Docker image**:

   ```bash
   docker pull sqlc/sqlc
   ```

2. **Initialize an SQLC configuration file**:

   ```bash
   docker run --rm -v "$(pwd):/src" -w /src sqlc/sqlc init
   ```

3. **Generate Go Code from SQL queries**:

   ```bash
   docker run --rm -v "$(pwd):/src" -w /src sqlc/sqlc generate
   ```

### Use Docker

Now that we made some api end points lets build an image for our backend

first lets create a docker file and pass next code:

    ```yaml
    FROM golang:1.20-alpine3.18

    WORKDIR /app

    COPY . .

    RUN go build -o main main.go

    EXPOSE 8080:8080

    CMD ["/app/main"]

    ```

To build the image run the next command

    ```bash
    docker build -t backend-masterclass:1.1 .
    ```

As you can see we giving a name for our docker image using the -t flag with a version also.

Now we have a bit of a problem the size of the image is a bit BIG, which is okay as long as our image works right no
the image contains alot of unnecessary things in it which are using to build the image but not needed after becouse golang
produces a binary file that can be run with our the need of nay language specific tools.

So now how can we make this less BIG ?

We have to use something called multi-stage builds, as the name suggests we have to build our image in multiple steps step one is download the golang image and the packages the our project needs and then build them, next step we copy the result binary file of our project and run it heres the code:

    ```yaml
    # step one
    FROM golang:1.20-alpine3.18 AS builder

    WORKDIR /app

    COPY . .

    RUN go build -o main main.go

    #step two
    FROM alpine

    WORKDIR /app

    COPY --from=builder /app/main .

    COPY app.env .

    EXPOSE 8080

    CMD ["/app/main"]
    ```

as you can see we build the app in the first step and all we do in the second step is just run the app

now lets create a container using our image:

    ```bash
        docker run --name backend-container -p 8080:8080 -e GIN_MODE=release 305c5c8595d5660950a525670746a8f7e1c77a34e614f62e4f83e03cce3c05e5
    ```

here we run our container with the name backend-container and exposing the ports using -p flag and define some variable with -e flag and the last long string is the id of our image
