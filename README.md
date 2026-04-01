# 🔍 FindIt – Lost & Found System

FindIt is a full-stack web application that helps people report lost and found items and reconnect them with their rightful owners.

---

## 🚀 Overview

Losing items is frustrating — and finding them again is often a matter of luck.  

**FindIt** simplifies this process by allowing users to:

- Report lost items  
- Report found items  
- Browse all reported items  
- Increase chances of reconnecting owners with their belongings  

---

## ✨ Features

- 🔐 User authentication (Sign up & Login)
- 📦 Create lost item reports
- 📍 Create found item reports
- 📋 View all reported items
- 🔎 Filter items (Lost / Found)
- 🌐 Clean and responsive UI

---

## 🛠️ Tech Stack

### Backend
- **Go (Golang)**
- **SQLite**
- **net/http** (standard library)

### Frontend
- **HTML**
- **CSS**
- **Vanilla JavaScript**

---

## 📁 Project Structure
findit/

│

├── internal/

│ ├── db/ # Database setup and migrations

│ ├── handlers/ # API handlers

│ ├── models/ # Data structures

│ └── routes/ # Routing logic

│

├── pkg/

│ └── utils/ # Helper utilities (JSON, JWT)

│

├── frontend/ # UI (HTML, CSS, JS)

│

├── main.go # Application entry point

└── findit.db # SQLite database (auto-created)


---

## ⚙️ Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/YOUR_USERNAME/findit.git
cd findit

go mod tidy

go run main.go

http://localhost:8080

```

### 2. Using Docker (Recommended)

#### Build and run with Docker Compose

```bash
# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

#### Build and run with Docker directly

```bash
# Build the Docker image
docker build -t findit-server .

# Run the container
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data --name findit-server findit-server

# View logs
docker logs -f findit-server

# Stop the container
docker stop findit-server
```

The application will be available at `http://localhost:8080`

API Endpoints
Authentication
| Method | Endpoint | Description       |
| ------ | -------- | ----------------- |
| POST   | /signup  | Register new user |
| POST   | /login   | Login user        |

Items
| Method | Endpoint | Description              |
| ------ | -------- | ------------------------ |
| POST   | /items   | Create item (lost/found) |
| GET    | /items   | Retrieve all items       |

Example Request

{

  "id": "1",

  "user_id": "u1",

  "type": "lost",

  "name": "Black Phone",

  "description": "iPhone with cracked screen",

  "location": "Nairobi",

  "date": "2026-04-01"

}

-Each item must have a unique ID

-The user_id must exist before creating an item

-The database is created automatically on first run