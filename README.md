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