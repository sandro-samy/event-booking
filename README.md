# Event Booking API

A small REST API for creating and managing events, built with Go, [Gin](https://gin-gonic.com/) and SQLite. Users register, log in to get a JWT, and can only see and change their own events.

## Getting started

Requires Go 1.27+.

1. Create a `.env` file in the project root:

   ```
   JWT_SECRET=some-long-random-string
   ```

2. Run the server:

   ```bash
   go run .
   ```

The API listens on `http://localhost:8080`. The SQLite database (`api.db`) is created automatically on first run.

To run it with Docker instead, see [DOCKER.md](DOCKER.md).

## Endpoints

| Method | Path             | Auth | Description                |
| ------ | ---------------- | ---- | -------------------------- |
| POST   | `/auth/register` | No   | Create an account          |
| POST   | `/auth/login`    | No   | Log in and get a JWT       |
| GET    | `/events`        | Yes  | List your events           |
| POST   | `/events`        | Yes  | Create an event            |
| GET    | `/events/:id`    | Yes  | Get one of your events     |
| PUT    | `/events/:id`    | Yes  | Update one of your events  |
| DELETE | `/events/:id`    | Yes  | Delete one of your events  |

Authenticated routes need the token from `/auth/login` in the header:

```
Authorization: Bearer <token>
```

Tokens expire after 24 hours. Accessing another user's event returns `403`.

### Example bodies

Register / login:

```json
{ "email": "user@example.com", "password": "Password123!" }
```

Create / update an event:

```json
{
  "name": "Go Meetup",
  "description": "Monthly meetup",
  "location": "Cairo",
  "date_time": "2026-06-01T10:00:00Z"
}
```

Ready-made requests are in [api-test/](api-test/) (for the VS Code REST Client or JetBrains HTTP Client).

## Project structure

```
main.go        entry point: loads .env, opens the DB, starts the server
db/            SQLite connection and table setup
models/        User and Event types and their queries
routes/        HTTP handlers and route registration
middleware/    JWT auth middleware
utils/         password hashing, JWT and Bearer header helpers
api-test/      .http files for trying the API by hand
```
