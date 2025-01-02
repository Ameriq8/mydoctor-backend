# Backend Server

## Prerequisites

Before running the server, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.20 or later)
- [Docker](https://www.docker.com/products/docker-desktop) (for running dependencies like PostgreSQL, Prometheus, etc.)
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
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=mydoctor
DB_SSLMODE=disable
DB_DRIVER=postgres
```

Copy the example `.env` file if it does not already exist:
```bash
cp .env.example .env
```

Edit the `.env` file to suit your environment.

### 3. Run the Server

Start the backend server:
```bash
go run cmd/app/main.go
```

The server should now be running at [http://localhost:8080](http://localhost:8080).

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

### 2. Verify Services

- **Backend**: [http://localhost:8080](http://localhost:8080)
- **Prometheus**: [http://localhost:9090](http://localhost:9090)
- **Alertmanager**: [http://localhost:9093](http://localhost:9093)

### 3. Stop Services

To stop the services:
```bash
docker-compose down
```

---

## Monitoring and Debugging

### Prometheus Metrics
Prometheus metrics are exposed at `/metrics`:
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
