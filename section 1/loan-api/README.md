# Loan Origination System API

This project implements a simple REST API for a loan origination system using Go and the Gin web framework. It provides endpoints for managing loan applications, including submission, retrieval, status updates, and document uploads. Data is stored in memory for demonstration purposes.

---

## Project Structure 

```
loan-api/
├── main.go                     # Main entry point of the application
├── handler/                    # Contains HTTP handler functions
│   └── loan.go                 # Handlers for loan application endpoints
├── usecase/                    # (Placeholder) For business logic that orchestrates store operations
├── middleware/                 # Custom Gin middleware functions
│   ├── auth.go                 # Authentication middleware
│   ├── logger.go               # Request logging middleware
│   └── error_handler.go        # Custom error recovery middleware
├── model/                      # Data structures/models
│   └── loan.go                 # LoanApplication struct and error response format
├── store/                      # Data storage layer
│   └── memory.go               # In-memory implementation of data storage
├── routes/                     # Defines API routes
│   └── routes.go               # Centralized route setup
└── tests/                      # Unit tests for the API
    └── loan_test.go            # Tests for loan application endpoints
```

---

## Setup Instructions

### 1. Clone & Setup

```bash
git clone https://github.com/nanda91/loan-api.git
cd loan-api
```

### 2. Run

```bash
go run main.go
```

### 3. Run Tests

```bash
go test ./tests... -v
```
---

## API Documentation

### Endpoints

| Method | Endpoint                              | Description                        |
|--------|---------------------------------------|------------------------------------|
| GET    | `/loan-applications`                  | List all applications              |
| GET    | `/loan-applications/:id`              | Get specific application           |
| POST   | `/loan-applications`                  | Submit new loan application        |
| PUT    | `/loan-applications/:id/status`       | Update loan status                 |
| POST   | `/loan-applications/:id/documents`    | Upload documents (multipart form)  |

---

1. List All Loan Applications
   - Endpoint: `GET /loan-applications`
   - Query Parameters:
       - page (optional, integer): Page number (default: 1)
       - limit (optional, integer): Number of applications per page (default: 10)
       - status (optional, string): Filter applications by status (e.g., pending, approved, rejected, under_review). Case-insensitive.
   - Authentication: Required
   - `200` OK: A JSON array of LoanApplication objects. SSN is masked.
        ```text
        [
          {
            "id": 1,
            "applicant_name": "Nanda",
            "applicant_ssn": "XXX-XX-6789",
            "loan_amount": 50000,
            "loan_purpose": "Home Renovation",
            "annual_income": 75000,
            "credit_score": 720,
            "status": "pending",
            "submitted_at": "2023-10-27T10:00:00Z",
            "documents_uploaded": []
          }
        ]
        ```
2. Get Specific Loan Application
   - Endpoint: `GET /loan-applications/{id}`
   - URL Parameters:
       - `id`(integer, required): The ID of the loan application.
   - Authentication: Required
   
   - `200` OK: A JSON array of LoanApplication objects. SSN is masked.
        ```text
          {
            "id": 1,
            "applicant_name": "Nanda",
            "applicant_ssn": "XXX-XX-6789",
            "loan_amount": 50000,
            "loan_purpose": "Home Renovation",
            "annual_income": 75000,
            "credit_score": 720,
            "status": "pending",
            "submitted_at": "2023-10-27T10:00:00Z",
            "documents_uploaded": []
          }
        ```
   - Error Responses
     - 400 Bad Request: If the ID is not a valid integer.
     - 401 Unauthorized: If authentication token is missing or invalid.
     - 404 Not Found: If no application with the given ID exists.
     - 500 Internal Server Error: For unexpected server errors


3. Submit New Loan Application
    - Endpoint: `POST /loan-applications`
   - Authentication: Required
   - Request Body
        ```text
        {
          "applicant_name": "Nanda",
          "applicant_ssn": "987-65-4321",
          "loan_amount": 75000.50,
          "loan_purpose": "Business Expansion",
          "annual_income": 120000.00,
          "credit_score": 800
        }
        ```
   - `201` Created: The newly created LoanApplication object. SSN is masked
        ```text
        [
          {
            "id": 1,
            "applicant_name": "Nanda",
            "applicant_ssn": "XXX-XX-6789",
            "loan_amount": 50000,
            "loan_purpose": "Home Renovation",
            "annual_income": 75000,
            "credit_score": 720,
            "status": "pending",
            "submitted_at": "2023-10-27T10:00:00Z",
            "documents_uploaded": []
          }
        ]
        ```
4. Update Application Status
    - Endpoint: `PUT /loan-applications/{id}/status`
   - Authentication: Required
   - URL Parameters
     - `id`(integer, required): The ID of the loan application.
   - Request Body
        ```text
        {
          "applicant_name": "Nanda",
          "applicant_ssn": "987-65-4321",
          "loan_amount": 75000.50,
          "loan_purpose": "Business Expansion",
          "annual_income": 120000.00,
          "credit_score": 800
        }
        ```
   - `200` OK: The updated LoanApplication object. SSN is masked
        ```text
        [
          {
            "id": 1,
            "applicant_name": "Nanda",
            "applicant_ssn": "XXX-XX-6789",
            "loan_amount": 50000,
            "loan_purpose": "Home Renovation",
            "annual_income": 75000,
            "credit_score": 720,
            "status": "pending",
            "submitted_at": "2023-10-27T10:00:00Z",
            "documents_uploaded": []
          }
        ]
        ```
5. Upload Supporting Documents
    - Endpoint: `POST /loan-applications/{id}/documents`
    - Authentication: Required
    - URL Parameters
        - `id`(integer, required): The ID of the loan application.
    - Request form-data
    ```text
    
      "document" [file]: "payslip.pdf"
    
    ```
    - `200` OK: The updated LoanApplication object. SSN is masked
         ```text
         [
           {
             "id": 1,
             "applicant_name": "Nanda",
             "applicant_ssn": "XXX-XX-6789",
             "loan_amount": 50000,
             "loan_purpose": "Home Renovation",
             "annual_income": 75000,
             "credit_score": 720,
             "status": "pending",
             "submitted_at": "2023-10-27T10:00:00Z",
             "documents_uploaded": [payslip.pdf]
           }
         ]
   


##  Middleware

-  **Auth Middleware** (simulated, accepts bearer token)
-  **Logger Middleware** for request logging
-  **Panic Recovery + Error Handler**
---

##  System Design Overview

### Diagram

```text
+----------------+       +------------------------------------+       +-----------------+
|                |       |                                    |       |                 |
|    Client      |------>|       Loan Origination API         |------>|   In-Memory     |
| (Web/Mobile App)|      |                                    |       |    Database     |
|                |       |                                    |       | (map[int]LoanApp)|
+----------------+       |  +------------------------------+  |       +-----------------+
                         |  |  Middleware:                 |  |
                         |  |  - Error Recovery            |  |
                         |  |  - Request Logging           |  |
                         |  |  - Authentication            |  |
                         |  +------------------------------+  |
                         |  +------------------------------+  |
                         |  |  API Endpoints (Handlers):   |  |
                         |  |  - GET /loan-applications    |  |
                         |  |  - GET /loan-applications/{id}|  |
                         |  |  - POST /loan-applications   |  |
                         |  |  - PUT /loan-applications/{id}/status |  |
                         |  |  - POST /loan-applications/{id}/documents |  |
                         |  +------------------------------+  |
                         +------------------------------------+
```

---
