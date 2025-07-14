# Loan Origination Microservice

---

## Project Structure

```
loan-api/
├── main.go                      # Entry point
├── cache/                       # In-memory cache for credit reports
│   └── cache.go
├── metrics/                     # Prometheus metrics
│   └── metrics.go
├── model/                       # Domain models
│   └── model.go
├── repository/                  # repository (Query Optimization)
│   └── loan_repository.go
├── service/                     # Business logic (credit & underwriting services)
│   ├── credit_service.go
│   └── underwriting_service.go
```

---

## Setup Instructions

### Install dependencies
```bash
go mod tidy
```

### Run the application
```bash
go run main.go
```

The application will:
- Simulate 10 loan application evaluations
- Call `CreditService` (mocked)
- Serve Prometheus metrics on `localhost:2112/metrics`

---

## Prometheus Metrics
Exposed at `http://localhost:2112/metrics`

- `underwriting_approved_total`
- `underwriting_rejected_total`
- `credit_score_request_duration_seconds`

---

##  Features

### 1. Service Communication
- Uses `UnderwritingService` ↔ `CreditService` interface
- Circuit breaker using `github.com/sony/gobreaker`

### 2. Caching
- SSN → Credit Report cached for 24 hours (in-memory)
- Fallback to stale cache on failure

### 3. Error Handling
- If `CreditService` fails and no cache exists, returns error
- Uses circuit breaker to avoid flooding the credit service

---

## Diagram

```
+--------------------+      +------------------+
| UnderwritingService|<--->|   CreditService   |
|  (with cache + CB) |      +------------------+
|       ^            |               ^
|       |            |               |
|       v            |               v
| Prometheus Metrics |     In-Memory Cache (24h)
+--------------------+
```

---

## Optimization Notes

### Query Optimization (Section 2.2)

DATABASE OPTIMIZATION EXPLANATION:

1. **Indexing Strategy**:
    - `idx_loan_applications_status`: Composite index on (status, last_updated DESC) for efficient status queries with ordering
    - `idx_loan_applications_applicant_hash`: Index on applicant_ssn_hash for duplicate checking
    - `idx_loan_applications_last_updated`: Index for general time-based queries

2. **Query Optimization**:
    - Added pagination with LIMIT/OFFSET
    - Used SSN hash instead of plain SSN for indexing (better for PII protection)
    - Removed SELECT * and specified exact columns needed
    - Added ORDER BY for consistent results

3. **PII Protection**:
    - Store hashed SSN for indexing purposes
    - Encrypt actual SSN separately
    - Use hash for all database operations

4. **Performance Improvements**:
    - Composite indexes prevent full table scans
    - Pagination reduces memory usage
    - Column selection reduces I/O
    - Proper ordering utilizes indexes

---


the solution provides:
- Efficient service communication with caching and circuit breaker
- Optimized database queries with proper indexing
- Comprehensive monitoring with Prometheus metrics
- Graceful error handling and fallback mechanisms
- Scalable architecture for high-volume processing (500+ requests/minute)