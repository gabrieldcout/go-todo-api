# 🧩 go-todo-api

A **RESTful API** written in **Go (Golang)** for managing tasks (**To-Do List**) with **JWT authentication**.  
The project includes authentication, token refreshing, a complete task CRUD, and integration with **Docker** and **CI/CD** for automatic unit testing.

---

## 🚀 Features

- ✅ **JWT Authentication**
  - User login returns a **JWT token**.
  - Some routes require authentication via token.
  - **Refresh token** endpoint to renew tokens without re-login.
  
- 📝 **Task Management (CRUD)**
  - `GET /tasks` — Lists all tasks.
  - `POST /tasks` — Creates a new task.
  - `PUT /tasks/{id}` — Updates an existing task.
  - `DELETE /tasks/{id}` — Deletes a task.

- ⚙️ **CI/CD**
  - Pipeline configured to automatically run **unit tests** on each commit or pull request.

- 🐳 **Docker and Docker Compose**
  - Includes `Dockerfile` and `docker-compose.yml` to easily spin up and test the application.

---

## 🧠 Technologies Used

- **Go** latest
- **JWT (JSON Web Token)** — authentication and authorization
- **Gin** — fast and lightweight HTTP framework
- **PostgreSQL / SQLite** — data persistence (depending on configuration)
- **Docker** & **Docker Compose**
- **GitHub Actions** (or another CI tool) — automatic test execution

---

## 🔑 Authentication

After login, the API returns a **JWT token** that must be included in the header of authenticated requests:

```http
Authorization: Bearer <your_token_here>
```

### Authentication Routes

| Method | Endpoint         | Description |
|--------|------------------|--------------|
| `POST` | `/login`         | Logs in and returns a JWT |
| `POST` | `/refresh`       | Generates a new token when the current one expires |

---

## 🧰 Main Routes

| Method | Endpoint      | Description | Authentication |
|--------|----------------|--------------|----------------|
| `GET`    | `/tasks`          | Lists all tasks | 🔒 Yes |
| `POST`   | `/tasks`          | Creates a new task | 🔒 Yes |
| `GET`    | `/tasks/{id}`     | Retrieves a specific task | 🔒 Yes |
| `PUT`    | `/tasks/{id}`     | Updates a task | 🔒 Yes |
| `DELETE` | `/tasks/{id}`     | Deletes a task | 🔒 Yes |

---

## 🧪 Testing

Unit tests are automatically executed via CI.  
To run them manually:

```bash
go test ./...
```

---

## 🐳 Running with Docker Compose

### 1️⃣ Start the application

```bash
docker-compose up --build
```

The application will be available at:  
👉 `http://localhost:8080`

---

## ⚙️ Database Connection (Client Access)

HOST: localhost  
USER: postgres  
NAME: tasks_db  
PORT: 5432  
PASSWORD: postgres

---

## 🧑‍💻 Contributing

1. Fork the project  
2. Create a branch (`git checkout -b feature/new-feature`)  
3. Commit your changes (`git commit -m 'feat: add new feature'`)  
4. Push to your branch (`git push origin feature/new-feature`)  
5. Open a **Pull Request**

---
