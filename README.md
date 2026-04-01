# FindIt – Lost & Found System

A backend system for reporting and recovering lost items.

## Tech Stack

- **Language:** Go
- **Framework:** net/http (standard library)
- **Database:** SQLite
- **Auth:** bcrypt password hashing

## Running

```bash
go run cmd/server/main.go
```

Server starts on `:8080`.

## API Endpoints

### Auth

| Method | Endpoint  | Description         |
|--------|-----------|---------------------|
| POST   | /signup   | Register a new user |
| POST   | /login    | Log in              |

### Items

| Method | Endpoint   | Description            |
|--------|------------|------------------------|
| POST   | /items     | Report a lost/found item |
| GET    | /items     | List all reported items  |

### Request Bodies

**POST /signup**
```json
{
  "id": "u1",
  "email": "user@example.com",
  "password": "1234"
}
```

**POST /login**
```json
{
  "email": "user@example.com",
  "password": "1234"
}
```

**POST /items**
```json
{
  "id": "i1",
  "title": "Black Wallet",
  "description": "Found near the bus stop",
  "type": "found",
  "user_id": "u1"
}
```

## Project Structure

```
findit/
├── cmd/server/main.go
├── internal/
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   └── item_handler.go
│   ├── models/
│   │   ├── user.go
│   │   └── item.go
│   ├── routes/routes.go
│   └── db/db.go
├── pkg/utils/response.go
├── go.mod
└── README.md
```
