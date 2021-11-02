> "Truth can only be found in one place: the code." - Robert C. Martin

# Data Warehouse
Data Warehouse is a simple backend application that exposes REST API and JWT-authenticated GraphQL API that allow for CRUD operations of the user entity. Data is stored in Postgres database and Redis. Shortly, more features will be added.
### Features in progress:
- Import of csv file with marketing data using Goroutines
- GraphQL API to query marketing data in a generic and efficient way

### Why Data Warehouse?

This is a learning exercise for Go, GraphQL, Redis, Docker, and related technologies.

### Tech
- Go v1.17
- GraphQL
- Postgres
- Redis
- Docker
- Gorm

# Getting Started

### Configurations:
- Project (server, database, cache) - config.yaml
- Docker (postgres, redis, pgadmin4) - docker/docker-compose.yaml

### Running the project locally
First, start docker environments (in `./docker`):
```
docker-compose up
```
Database schema migration is executed automatically via gorm.io.

Next, run the project (by default `config.yaml` is loaded):
```
go run ./cmd/server.go
```
Or you can load other config file:
```
go run ./cmd/server.go -config config.yaml
```
