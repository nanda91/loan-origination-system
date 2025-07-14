#  Loan Document Processor

---

##  Project Structure

```
loan-doc-processor/
├── main.go                       # App entrypoint
├── handler/
│   └── loat.go                   # Router and endpoint handler
├── model/
│   └── document.go               # Structs for jobs and results
├── processor/
│   └── document_processor.go     # Worker logic and document parsing
├── queue/
│   └── job_queue.go              # Job queue using Go channels
├── utils/
│   └── pdf_utils.go              # PDF extraction utilities
├── go.mod                        # Dependencies
└── README.md                     # Project documentation
```

---

## Installation

1. Install dependencies and run the app:

```bash
go mod tidy
go run main.go
```

---

## API Endpoint

### `POST /loan-applications/:id/documents`

Accepts a PDF file for a specific loan application.

#### Headers:

* `Content-Type: multipart/form-data`

#### Form Data:

| Field           | Type   | Required | Description                                |
| --------------- | ------ | -------- | ------------------------------------------ |
| `file`          | File   | ✅        | The PDF document to upload                 |
| `document_type` | String | ✅        | `bank_statement`, `pay_stub`, `tax_return` |

#### Example `curl`:

```bash
curl -X POST http://localhost:8080/loan-applications/123/documents \
  -F "file=@example.pdf" \
  -F "document_type=bank_statement"
```

#### Response:

```json
{
  "status": "queued",
  "file": "uploaded_filename.pdf"
}
```

---

## How It Works

1. File is uploaded and stored temporarily.
2. A `DocumentJob` is pushed to a buffered channel.
3. Background workers (3 by default) pull jobs and:

    * Extract plain text using PDF parser
    * Simulate extracting important data (for now: content snippet)
    * Call the job's callback function (e.g., logging result)

---