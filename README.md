# Starting a Golang backend service :

## Bank Backend API

A simple backend API for managing bank accounts, user authentication, and money transfers.

## Features

- User registration and authentication using JWT.
- Creating and managing bank accounts.
- Secure money transfers between accounts.


### Prerequisites

- Go (version 1.21.5)
- PostgreSQL (16.1)

## Project Includes the following : 
- [x] Building a simple API (Creating Bank accounts) .
- [x] Postgres
- [x] JWT 
- [x] Docker
- [x] Specific Architecture / " Design Pattern  " 
- [x] Refactoring
- [x] Login 
- [x] registration 
- [x] Tests
- [x] Dockerizing Go " Add Dockerfile "  

___
### Useful Commands : 
_Creating Docker Network for our service :_
```bash
docker network create backend_go_network
```
_Initializing Postgres Docker container in our created Network :_
```bash
docker run --name backend-go-postgres --network backend_go_network -e POSTGRES_PASSWORD=backend -p 5432:5432 -d postgres
```
_Running containers along with their associated networks_

```bash
docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Networks}}"
```

_Build the go Docker Image :_

```bash
docker build -t backend_go .
```

_Then you good to go and run the container :_
```bash
docker run -p 3000:3000 backend_go
```

___