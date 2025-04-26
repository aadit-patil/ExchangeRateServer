# EXCHANGE RATE SERVER

A scalable Go microservice for real-time and historical currency conversion with:

- In-memory caching
- MySQL persistence
- External API fallback
- Prometheus monitoring

## Prerequisites

- Go 1.22+
- Docker
- Docker Compose

## Getting Started

### 1. Clone the repository

```bash
git clone [https://github.com/aadit-patil/ExchangeRateServer.git](https://github.com/aadit-patil/ExchangeRateServer.git)
cd ExchangeRateServer
```

### 2. Start the services
```bash
docker-compose up --build
```

This will start:

-   `db` (MySQL 8 container)
-   `app` (ExchangeRateServer running at port `8088`)

### API Usage

#### Endpoint URL
```bash
GET /convert
```
### Query Parameters
| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| from | string | Yes | Base currency code (e.g., USD) |
| to | string | Yes | Target currency code (e.g., INR) |
| amount | float | No | Amount to convert (optional) |
| date | string | No | Date in YYYY-MM-DD format (optional, defaults to today) |

## Example API Calls

### 1. Get today's conversion rate

`curl --location 'http://localhost:8088/convert?from=USD&to=INR'`

### 2. Convert specific amount

`curl --location 'http://localhost:8088/convert?from=USD&to=INR&amount=1000'`

### 3. Conversion for a historical date

`curl --location 'http://localhost:8088/convert?from=USD&to=INR&amount=1000&date=2024-04-25'`

### 4. Get only historical rate (without amount)

`curl --location 'http://localhost:8088/convert?from=USD&to=EUR&date=2024-04-25'`

## Prometheus Metrics

Metrics are available at:

`http://localhost:8088/metrics`

### Exposed Metrics

-   cache_hits_total
-   cache_misses_total
-   db_queries_total
-   external_api_requests_total