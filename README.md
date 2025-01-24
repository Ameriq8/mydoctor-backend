# Backend Server

## Prerequisites

Before running the server, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.20 or later)
- [Docker](https://www.docker.com/products/docker-desktop) (for running dependencies like PostgreSQL, Redis, Prometheus, etc.)
- [Artillery](https://www.artillery.io/docs) (for performance testing)

---

## How to Run

### 1. Install Dependencies

Ensure all Go module dependencies are installed:
```bash
go mod tidy
```

### 2. Set Up Environment Variables

Ensure your `.env` file is properly configured with the necessary environment variables. An example is provided below:

```
GIN_MODE=debug
SERVER_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Ameriq81
DB_NAME=mydoctor
DB_SSLMODE=disable
DB_DRIVER=postgres
REDIS_HOST=redis
REDIS_PORT=6379
```

Copy the example `.env` file if it does not already exist:
```bash
cp .env.example .env
```

Edit the `.env` file to suit your environment.

### 3. Run the Server

Start the backend server with Docker:
```bash
docker-compose up --build
```

This will start all necessary services including the backend, PostgreSQL, Redis, Prometheus, and Alertmanager. The server should now be running at [http://localhost:8080](http://localhost:8080).

---

## API Endpoints

### Test Endpoint
- **Endpoint**: `/ping`
- **Method**: `GET`
- **Description**: Returns a test response to verify the server is running.

Example request:
```bash
curl http://localhost:8080/ping
```

Expected response:
```json
{
  "message": "pong",
}
```

---

## Performance Testing

### Prerequisites

Ensure Artillery is installed globally:
```bash
npm install -g artillery
```

### Run Performance Tests

Execute load tests using Artillery:
```bash
artillery run artillery-config.yml --record --key your-artillery-key
```

- Replace `your-artillery-key` with your Artillery API key if you are recording results.
- The `artillery-config.yml` file should define your performance testing scenarios.

### Sample `artillery-config.yml`

```yaml
config:
  target: "http://localhost:8080" # Base URL of the application
  phases:
    - name: "Ramp-up to Moderate Load"
      duration: 30
      arrivalRate: 20 # Ramp-up to 20 requests per second
    - name: "Steady Load"
      duration: 60
      arrivalRate: 50 # Steady load of 50 requests per second
    - name: "High Load Test"
      duration: 30
      arrivalRate: 200 # High load of 200 requests per second

  defaults:
    headers:
      Content-Type: "application/json" # Set content type if needed
    http:
      timeout: 5000 # Timeout in ms
      retry: 2       # Retry up to 2 times if a request fails

scenarios:
  - name: "Ping API Test"
    flow:
      - get:
          url: "/ping" # Endpoint to test
```

---

## Running in Docker

### 1. Build and Start Services

Use `docker-compose` to build and start all services:
```bash
docker-compose up --build
```

This command will build the containers and bring them online, including the following services:

- **Backend**: [http://localhost:8080](http://localhost:8080)
- **PostgreSQL**: Exposed on port 5432 (inside Docker network)
- **Redis**: Exposed on port 6379 (inside Docker network)
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Alertmanager**: [http://localhost:9093](http://localhost:9093)
- **SMTP**: Exposed on port 1025 (for email testing)

### 2. Verify Services

Verify that all the services are up and running:

- Backend: [http://localhost:8080](http://localhost:8080)
- Prometheus: [http://localhost:9090](http://localhost:9090)
- Alertmanager: [http://localhost:9093](http://localhost:9093)

### 3. Stop Services

To stop the services and clean up:
```bash
docker-compose down
```

This will stop and remove all containers, but retain the volume data (such as your PostgreSQL data).

---

## Monitoring and Debugging

### Prometheus Metrics
Prometheus metrics are exposed at `/metrics` on the backend service:
- Example: [http://localhost:8080/metrics](http://localhost:8080/metrics)

### Debugging Tips

- **View Logs**:
  ```bash
  docker logs gin_app
  ```

- **Inspect Running Containers**:
  ```bash
  docker ps
  ```

- **Access Docker Shell**:
  ```bash
  docker exec -it gin_app sh
  ```

---

## Unit and Integration Testing

### Unit Tests

To run unit tests, ensure that your database is not required (or mocked) since unit tests are isolated from external dependencies.

#### 1. Install Test Dependencies

Make sure you have installed all necessary test dependencies:
```bash
go mod tidy
```

#### 2. Running Unit Tests

To run unit tests for your Go project:
```bash
go test ./... -v
```

This will execute the unit tests in all your packages and show detailed output for each test case.

### Integration Tests

For integration tests that require a real database or external dependencies (PostgreSQL and Redis in this case), you can use Docker to spin up the required services.

#### 1. Set Up Docker for Integration Tests

Make sure you have the necessary services up and running in Docker. You can use the `docker-compose.yml` file to bring up PostgreSQL, Redis, and any other required services.

```yaml
version: "3.8"
services:
  app:
    container_name: gin_app
    restart: always
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=debug
      - SERVER_PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=Ameriq81
      - DB_NAME=mydoctor
      - DB_SSLMODE=disable
      - DB_DRIVER=postgres
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    networks:
      - monitoring
    depends_on:
      - postgres
      - prometheus
      - redis

  postgres:
    container_name: postgres
    image: postgres:15
    restart: always
    ports:
      - "3002:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: Ameriq81
      POSTGRES_DB: mydoctor
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - monitoring

  redis:
    container_name: redis
    image: redis:6.2
    restart: always
    ports:
      - "6379:6379"
    networks:
      - monitoring

  prometheus:
    container_name: prometheus
    restart: always
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
    networks:
      - monitoring

  alertmanager:
    container_name: alertmanager
    image: prom/alertmanager
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml
    ports:
      - "9093:9093"
    networks:
      - monitoring

volumes:
  postgres_data:
  
networks:
  monitoring:
    driver: bridge
```

#### 2. Running Integration Tests

To run your integration tests that interact with the database:
```bash
go test -tags=integration ./... -v
```

Make sure your integration tests are tagged correctly, or you can use the `-tags` flag to define the correct tags for integration tests.

#### 3. Example Integration Test

Below is an example of an integration test for creating a user in the database.

```go
package repositories_test

import (
    "server/internal/models"
    "server/internal/repositories"
    "testing"

    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

var db *sqlx.DB

func setupDB(t *testing.T) {
    dsn := "host=localhost port=5432 user=postgres password=Ameriq81 dbname=mydoctor sslmode=disable"
    var err error
    db, err = sqlx.Connect("postgres", dsn)
    require.NoError(t, err)
}

func TestCreateUserIntegration(t *testing.T) {
    setupDB(t)

    repo := repositories.NewAuthRepository(db)

    user := &models.User{
        Name:  "John Doe",
        Email: "john.doe@example.com",
        Image: "image.png",
    }

    createdUser, err := repo.CreateUser

(user)
    require.NoError(t, err)

    assert.Equal(t, "John Doe", createdUser.Name)
    assert.Equal(t, "john.doe@example.com", createdUser.Email)
}
```

This test checks if the user creation functionality works by interacting with the real database running in Docker.
