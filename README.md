# Microservices Project with Krakend API Gateway

---

## Project Structure

```
.
├── krakend/
│   └── krakend.json          # API Gateway configuration
├── listing/                  # Listing service (PostgreSQL)
│   ├── Dockerfile
│   └── main.go
├── user/                     # User service (MongoDB)
│   ├── Dockerfile
│   └── main.go
└── docker-compose.yml
```

---

## Prerequisites

* Docker & Docker Compose
* Go 1.24+ (for building services)
* MongoDB & PostgreSQL (via Docker)

---

## Environment Variables

### User Service (MongoDB)

* `LISTEN_PORT` : port the service listens on (e.g., `:8889`)
* `MONGO_USER` : MongoDB username (e.g., `root`)
* `MONGO_PASSWORD` : MongoDB password (e.g., `12345`)
* `MONGO_ADDR` : MongoDB host (service name in Docker, e.g., `mongo`)
* `MONGO_PORT` : MongoDB port (e.g., `27017`)
* `MONGO_DATABASE` : Database name (e.g., `mydb`)

### Listing Service (PostgreSQL)

* `LISTEN_PORT` : port the service listens on (e.g., `:8888`)
* `POSTGRES_ADDR` : PostgreSQL host (service name in Docker, e.g., `postgres`)
* `POSTGRES_PORT` : PostgreSQL port (e.g., `5432`)
* `POSTGRES_DATABASE` : Database name (e.g., `mydb`)
* `POSTGRES_PASSWORD` : PostgreSQL password

---

## Docker Compose

Start all services including Krakend:

```bash
docker compose up -d
```

* **User Service:** [http://localhost:8889](http://localhost:8889)
* **Listing Service:** [http://localhost:8888](http://localhost:8888)
* **Krakend API Gateway:** [http://localhost:8080](http://localhost:8080)

---

## API Endpoints

### User Service

* `GET /users` → list all users
* `GET /users/{id}` → get user by ID
* `POST /users` → create a new user

### Listing Service

* `GET /listings` → list all listings
* `POST /listings` → create a new listing

### API Gateway (Krakend)

* `GET /users` → fetch from User Service
* `GET /users/{id}` → fetch user by ID from User Service
* `POST /users` → forward to User Service
* `GET /listings` → fetch from Listing Service
* `POST /listings` → forward to Listing Service

---

## Notes

1. Make sure the service names in Docker Compose (`mongo`, `postgres`, `user`, `listing`) match the hostnames used in your environment variables and Krakend config.
2. Krakend acts as a reverse proxy and forwards requests to microservices.
3. MongoDB and PostgreSQL are configured using Docker volumes for persistence.

---

## Commands

* Build all services:

```bash
docker compose build
```

* Stop and remove containers:

```bash
docker compose down
```

* View logs:

```bash
docker compose logs -f
```