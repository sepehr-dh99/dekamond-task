# Dekamond Task – OTP-Based Authentication Service (Golang 1.24.5)

This project is a backend service implemented in **Golang 1.24.5** for **OTP-based login and registration**, along with **basic user management**.  
It uses **in-memory storage** for simplicity, supports **JWT authentication**, and enforces **rate limiting** on OTP requests.

---

## **Features**

- **OTP Login & Registration**
  - Users request OTP by phone (Iran format: `09XXXXXXXXX`)
  - OTP valid for **2 minutes**, printed to **console**
  - Auto-registers new users, logs in existing ones
  - Returns **JWT** upon successful OTP verification
- **Rate Limiting**
  - Max **3 OTP requests per phone** within **10 minutes**
- **User Management**
  - Retrieve **single user details**
  - Retrieve **paginated & searchable** user list
- **Swagger/OpenAPI**
  - Documented REST APIs with annotations
- **Dockerized**
  - Lightweight image with **Go 1.24.5**
  - Ready to run via `docker`

---

## **Project Structure**

```
deKamond-task/
├── main.go
├── controller/
│   ├── dto/
│   │   ├── auth.go
│   │   └── user.go
│   ├── auth.go
│   └── user.go
├── service/
│   └── user.go
├── docs/
│   ├── swagger.json
│   ├── swagger.yml
│   └── docs.go
├── middleware/
│   └── auth.go
├── model/
│   └── user.go
├── package/
│   ├── jwt/
│   │   └── jwt.go
│   ├── otp/
│   │   └── otp.go
│   ├── response/
│   │   └── response.go
│   ├── validator/
│   │   └── validator.go
│   └── rate_limiter/
│       └── rate_limiter.go
├── go.mod
├── go.sum
├── README.md
├── License
└── Dockerfile
```

---

## **Getting Started**

### **Prerequisites**

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/)

---

### **Run Locally**

```bash
go mod tidy
go run main.go
```

Server runs on: **http://localhost:8080**

---

### **Run with Docker**

```bash
docker build -t dekamond-task .
docker run -d -p 8080:8080 --name dekamond-task dekamond-task
```

---

## **API Endpoints**

### **1. Request OTP**

```bash
curl -X POST http://localhost:8080/auth/request-otp \
  -H "Content-Type: application/json" \
  -d '{"phone": "09123456789"}'
```

**Response (200 OK)**:

```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": null
}
```

_(Check server logs for OTP)_

**Response (429 Too Many Requests)**:

```json
{
  "success": false,
  "message": "too many requests"
}
```

---

### **2. Verify OTP (Login/Register)**

```bash
curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"phone": "09123456789", "otp": "123456"}'
```

**Response (200 OK)**:

```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "<JWT_TOKEN>"
  }
}
```

**Response (400/401)**:

```json
{
  "success": false,
  "message": "phone and otp required"
}
```

---

### **3. Get Single User Details**

```bash
curl -X GET http://localhost:8080/users/09123456789 \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json"
```

**Response (200 OK)**:

```json
{
  "success": true,
  "message": "User fetched successfully",
  "data": {
    "phone": "09123456789",
    "registered_at": "2025-08-24T17:00:00Z"
  }
}
```

**Response (404 Not Found)**:

```json
{
  "success": false,
  "message": "User not found"
}
```

---

### **4. Get Paginated & Searchable User List**

```bash
curl -X GET "http://localhost:8080/users?page=1&size=5&search=091" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Response (200 OK)**:

```json
{
  "success": true,
  "message": "Users fetched successfully",
  "data": {
    "total": 12,
    "page": 1,
    "size": 5,
    "users": [
      {
        "phone": "09123456789",
        "registered_at": "2025-08-24T17:00:00Z"
      }
    ]
  }
}
```

**Response (401 Unauthorized)**:

```json
{
  "success": false,
  "message": "missing token"
}
```

---

## **Swagger/OpenAPI**

Swagger docs are generated via [swaggo/swag](https://github.com/swaggo/swag):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

Docs available at `/swagger/index.html` (once integrated).  
**Note:** use `Bearer <token>` in Swagger UI Authorization dialog.

---

## **Why In-Memory Storage?**

- **Fast & simple** for demos and interviews.
- **No external DB setup** (saves time, simpler Dockerization).
- Can be swapped for Redis/Postgres in production.

---

## **License**

Apache 2.0 – Use freely for learning and demonstration.
