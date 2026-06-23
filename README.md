# mini-service

A small REST API in Go: JWT authentication and a posts CRUD with per-user ownership.
Built to practice backend fundamentals — auth, password hashing, middleware, and an ORM.

## Stack

- **Go** + [Gin](https://github.com/gin-gonic/gin) — HTTP router
- **GORM** — ORM / database layer
- **golang-jwt/jwt v5** — JWT access tokens
- **bcrypt** — password hashing

## Run locally

```bash
git clone https://github.com/mxalle/mini-service.git
cd mini-service

export JWT_SECRET="your-secret-here"   # required, used to sign tokens

go mod download
go run .
```

Server starts on `http://localhost:8000`.

## Environment

| Variable     | Description                          | Required |
| ------------ | ------------------------------------ | -------- |
| `JWT_SECRET` | Secret key for signing JWT tokens    | yes      |

## Endpoints

| Method   | Path          | Auth | Description                          |
| -------- | ------------- | ---- | ------------------------------------ |
| `POST`   | `/signup`     | no   | Register a new user                  |
| `POST`   | `/login`      | no   | Log in, returns a JWT                |
| `GET`    | `/posts`      | no   | List all posts                       |
| `POST`   | `/posts`      | yes  | Create a post                        |
| `DELETE` | `/posts/:id`  | yes  | Delete own post (403 if not owner)   |

Protected routes expect the token in the header:
