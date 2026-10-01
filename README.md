# Go URL Shortener

A simple URL shortener API built with Go and Chi.

The application generates unique short codes for long URLs, stores them in memory, redirects users to the original URL, and tracks how many times each shortened link has been visited.

## Features

* Create shortened URLs
* Generate random 6-character short codes
* Detect and retry short-code collisions
* Redirect shortened URLs to their original destination
* Track link clicks
* Retrieve link information and statistics
* Thread-safe in-memory storage using `sync.RWMutex`
* JSON API

## Tech Stack

* Go
* Chi Router
* `net/http`
* JSON
* In-memory storage
* `sync.RWMutex`

## Project Structure

```text
url-shortener/
└── main.go
```

Everything is currently contained in `main.go` to keep the project simple and focused.

## Getting Started

### 1. Clone the repository

```bash
git clone <your-repository-url>
cd url-shortener
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

```bash
go run .
```

The server runs on:

```text
http://localhost:8080
```

---

## API Endpoints

| Method | Endpoint        | Description                               |
| ------ | --------------- | ----------------------------------------- |
| `POST` | `/links`        | Create a shortened URL                    |
| `GET`  | `/link/{code}`  | Redirect to the original URL              |
| `GET`  | `/links/{code}` | Get link information and click statistics |

---

## Create a Short URL

### Request

```http
POST /links
Content-Type: application/json
```

Body:

```json
{
  "url": "https://go.dev"
}
```

### Response

```json
{
  "code": "x7Yq2p",
  "url": "https://go.dev",
  "clicks": 0
}
```

The short code is generated automatically.

---

## Redirect

Use the generated code:

```http
GET /link/x7Yq2p
```

The server responds with an HTTP `302` redirect to:

```text
https://go.dev
```

Every successful redirect increments the link's click count.

---

## Link Statistics

To retrieve information about a shortened URL:

```http
GET /links/x7Yq2p
```

Response:

```json
{
  "code": "x7Yq2p",
  "url": "https://go.dev",
  "clicks": 3
}
```

---

## How Short Codes Work

When a new URL is submitted, the application generates a random 6-character code using `crypto/rand`.

The code is generated from:

```text
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789
```

If the generated code already exists, the application generates another code and checks again.

The basic flow is:

```text
URL submitted
     ↓
Generate 6-character code
     ↓
Check if code already exists
     ↓
 ┌───────────────┐
 │ Code exists?  │
 └───────┬───────┘
         │
    Yes  │  No
     ↓   │   ↓
 Generate│  Save link
 another │
   code  │
         ↓
       Return
```

This collision-handling pattern is useful beyond URL shorteners. The same idea can be applied to invite codes, usernames, idempotency keys, and other unique identifiers.

## Concurrency

The link store uses `sync.RWMutex` to protect the shared map.

Read operations use:

```go
RLock()
```

while operations that modify the store use:

```go
Lock()
```

This prevents concurrent requests from accessing or modifying the map unsafely.

## Storage

Links are currently stored in memory.

That means all links are lost whenever the application restarts.

A production version would use persistent storage such as PostgreSQL, MongoDB, or another database.

## Future Improvements

Possible improvements include:

* Persistent database storage
* Link expiration
* Custom short codes
* Rate limiting
* URL validation
* Authentication
* A small frontend form
* More detailed analytics
* Automated tests
* Docker support

## What This Project Demonstrates

This project demonstrates several backend concepts:

* Building an HTTP API with Go
* Routing with Chi
* JSON request and response handling
* HTTP redirects
* Random identifier generation
* Collision handling
* Shared state and mutexes
* Pointer-based state mutation
* Basic API design
* In-memory data storage

## Running the Project

```bash
go run .
```

Then create a link:

```bash
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d "{\"url\":\"https://go.dev\"}"
```

Copy the generated code and visit:

```text
http://localhost:8080/link/<code>
```

To check statistics:

```text
http://localhost:8080/links/<code>
```

---

## Project Status

This project is intentionally small and focused. Its main goal is to practice Go backend fundamentals, particularly unique identifier generation, collision handling, HTTP redirects, shared state, and API design.
