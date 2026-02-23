# NexBoard Task Management API

A robust, scalable RESTful backend API for managing projects and tasks, heavily inspired by modern collaborative workspace tools like Trello, Jira, and Asana. It is built in Golang using the Gin framework and PostgreSQL for relational integrity.

Every entity uses a **Dual-ID Strategy** (`external_id` for secure external API rendering + `id` for fast internal relational mappings) to enforce maximum security, abstraction, and efficiency.

## 📋 Features

- JWT-based authentication
- User registration and login flows
- Multi-tenancy Workspace management
- Granular Role-Based Access Control (RBAC: owner, admin, member)
- Board and global Status definition (Master Data) capabilities  
- Full Task assignment mapping
- Automatic soft delete/cascade cleanup protections
- Dual-ID database separation implementation
- PostgreSQL database with automated migrations via sql-migrate
- Clean architecture layout (handlers, services, repositories, middleware)
- Native Railway CI/CD Nixpacks deployment compatibility

## 🛠️ Technology Stack

- **Language:** Go (Golang) 1.22+
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL
- **Database Driver:** [lib/pq](https://github.com/lib/pq)
- **Migrations:** [sql-migrate](https://github.com/rubenv/sql-migrate)
- **Authentication:** JWT (JSON Web Token)
- **Encryption Tools:** bcrypt, UUIDv4
- **Environment Config:** godotenv

## � Project Structure

```
├── main.go                 # Application entry point/router definition
├── config/
│   ├── config.go          # Configuration loader
│   └── database.go        # Database connection pool wrapper
├── dbconfig.yml            # SQL-migrate binding mapper
├── railway.toml            # Automatic remote deployment configurations
├── models/
│   ├── user.go           # Profile bindings schema
│   ├── workspace.go      # Projects group abstraction
│   ├── workspace_member.go # RBAC participant mapping
│   ├── board.go          # Workspace subdivisions
│   ├── status.go         # Global column state trackers
│   └── task.go           # Base unit items schema
├── handlers/
│   ├── auth_handler.go   
│   ├── workspace_handler.go 
│   ├── board_handler.go   
│   ├── status_handler.go   
│   └── task_handler.go   
├── middleware/
│   └── auth_jwt.go        # JWT Context validation middleware
├── repositories/
│   ├── user_repository.go     
│   ├── workspace_repository.go 
│   ├── board_repository.go     
│   ├── status_repository.go     
│   └── task_repository.go     
├── services/
│   ├── auth_service.go        
│   ├── workspace_service.go    
│   ├── board_service.go        
│   ├── status_service.go        
│   └── task_service.go        
├── utils/
│   ├── jwt.go            
│   ├── password.go       
│   ├── response.go       
│   └── uuid.go       
└── migrations/
    ├── 001_create_users.sql
    ├── 002_create_workspaces.sql
    ├── 003_create_workspace_members.sql
    ├── 004_create_boards.sql
    ├── 005_create_statuses.sql
    └── 006_create_tasks.sql
```

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/your-username/nexboard-be.git
   cd nexboard-be
   ```

2. **Setup PostgreSQL database**

   Ensure PostgreSQL is running locally, then create a target database:
   ```bash
   # Create database
   createdb nexboard

   # Or using psql
   psql -U postgres
   CREATE DATABASE nexboard;
   ```

3. **Configure environment variables**

   Copy `.env.example` to `.env` and update with your database credentials:

   ```bash
   cp .env.example .env
   ```

   Edit `.env`:

   ```env
   PORT=8080
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/nexboard?sslmode=disable
   JWT_SECRET=your-super-secret-jwt-key
   ```

4. **Install Tools & Dependencies**

   ```bash
   go mod download
   go install github.com/rubenv/sql-migrate/...@latest
   ```

5. **Run database migrations**

   ```bash
   sql-migrate up -env=development
   ```

6. **Run the application**
   ```bash
   go run main.go
   ```

The server will start at `http://localhost:8080/```

## 🗺️ Entity Relationship Diagram (ERD)
![ERD](ERD_NexBoard.png)

## �️ Deployment to Railway

1. Push your code into a **GitHub repository**
2. On Railway, instantiate a new Project: **Deploy from GitHub repository** -> Your Repo.
3. Add a **PostgreSQL Database** within the same Railway project environment.
4. Go to your API project **Variables** tab (Service settings) in Railway and insert the exact keys found in `.env` bindings.

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

### Authentication

Almost all endpoints (except `/api/users/login` and `/api/users/register`) require a JWT token injected explicitly inside the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

---

### 🔐 Authentication Endpoints

#### 1. User Registration 

```http
POST /api/users/register
Content-Type: application/json

{
  "name": "Alice Developer",
  "email": "alice@test.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "external_id": "a1b2c3d4-e5f6-g7h8",
  "name": "Alice Developer",
  "email": "alice@test.com",
  "created_at": "2026-02-15T10:00:00Z"
}
```

#### 2. User Login

```http
POST /api/users/login
Content-Type: application/json

{
  "email": "alice@test.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "external_id": "a1b2c3d4-e5f6-g7h8",
    "name": "Alice Developer",
    "email": "alice@test.com"
  }
}
```

#### 3. Current User Profile

```http
GET /api/users/profile
Authorization: Bearer <token>
```

---

### 📂 Workspace Endpoints

#### 1. Create Workspace
_Automatically adds the invoking User as the 'owner'_

```http
POST /api/workspaces
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Engineering Team"
}
```

**Response (201 Created):**
```json
{
  "external_id": "w9x8y7z6",
  "name": "Engineering Team",
  "owner_external_id": "a1b2c3d4-e5f6-g7h8"
}
```

#### 2. Get Workspaces

```http
GET /api/workspaces
Authorization: Bearer <token>
```

#### 3. Update Workspace Details (Owner only)

```http
PUT /api/workspaces/w9x8y7z6
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Global Engineering Force"
}
```

#### 4. Delete Workspace 

```http
DELETE /api/workspaces/w9x8y7z6
Authorization: Bearer <token>
```

---

### 👥 Workspace Member Endpoints

#### 1. Invite Member

```http
POST /api/workspaces/w9x8y7z6/members
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_external_id": "b2c3d4a1",
  "role": "admin"
}
```

#### 2. Update Member Role
_Requires Owner level permissions constraints binding._

```http
PUT /api/workspaces/w9x8y7z6/members/b2c3d4a1
Authorization: Bearer <token>
Content-Type: application/json

{
  "role": "member"
}
```

#### 3. Remove Member
```http
DELETE /api/workspaces/w9x8y7z6/members/b2c3d4a1
Authorization: Bearer <token>
```

---

### 📋 Boards Endpoints

#### 1. Create Board

```http
POST /api/workspaces/w9x8y7z6/boards
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Sprint 1 Beta"
}
```

**Response (201 Created):**
```json
{
  "external_id": "b1b2b3b4",
  "workspace_external_id": "w9x8y7z6",
  "name": "Sprint 1 Beta"
}
```

#### 2. Get Workspace Boards

```http
GET /api/workspaces/w9x8y7z6/boards
Authorization: Bearer <token>
```

#### 3. Update/Delete Board
`PUT /api/boards/b1b2b3b4`
`DELETE /api/boards/b1b2b3b4`

---

### 🚥 Status Definitions (Global Metadata)

_Note: Trello style columns where the definitions globally span mapping available choices entirely._

#### 1. Create Status Configuration Label
```http
POST /api/statuses
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Done",
  "color": "#00ff00",
  "position": 3
}
```

#### 2. View Active Configurator Choices 
```http
GET /api/statuses
Authorization: Bearer <token>
```

---

### 📝 Tasks Endpoints

#### 1. Create Task

```http
POST /api/boards/b1b2b3b4/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Refactor router core",
  "priority": "high",
  "status_external_id": "s1s2s3s4",
  "assigned_to_external_id": "b2c3d4a1"
}
```

**Response (201 Created):**
```json
{
  "external_id": "t1t2t3t4",
  "title": "Refactor router core",
  "priority": "high"
}
```

#### 2. List Board Tasks Groupings 
_Automatically joins full Status / Assignee internal relations safely outward._

```http
GET /api/boards/b1b2b3b4/tasks
Authorization: Bearer <token>
```

**Response Context Preview:**
```json
[
  {
    "external_id": "t1t2t3t4",
    "title": "Refactor router core",
    "status": {
      "external_id": "s1s2s3s4",
      "name": "Done"
    },
    "assigned_to": {
      "external_id": "b2c3d4a1",
      "name": "Bob Programmer"
    }
  }
]
```

#### 3. Quick Move Options 
_Avoids passing arbitrary heavy bodies to fast lane updates._

```http
PATCH /api/tasks/t1t2t3t4/status
Content-Type: application/json

{
  "status_external_id": "s9s8s7s6"
}
```

```http
PATCH /api/tasks/t1t2t3t4/assign
Content-Type: application/json

{
  "assigned_to_external_id": null
}
```

---

