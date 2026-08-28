# short links

## Stack:

- golang + chi
- bcrypt (default salt)
- postgresql
- goose (migrations)
- redis
- kafka
- jwt authentication
- docker + docker compose
- git

Protected endpoints use /api and require JWT unless explicitly stated otherwise.
GET /r/{code} is public.

## Architecture

- `migrations/` - goose migrations
- `cmd/api/main.go` - start application
- `internal/<service_name>/` - users or links
- `internal/<service_name>/handler/` - HTTP handlers
- `internal/<service_name>/service/` - business logic, validation, Kafka events, cache logic
- `internal/<service_name>/repository/` - work with PostgreSQL
- `internal/<service_name>/cache/` - work with Redis
- `internal/<service_name>/events/` - work with Kafka
- `internal/<service_name>/dto/` - DTOs for requests and responses
- `Dockerfile` - build application image
- `docker-compose.yaml` - run backend, PostgreSQL, Redis and Kafka

Dockerfile - build an image
docker-compose.yaml - up kafka, redis, backend together

## Database scheme

### links

| Column          | Description               |
| --------------- | ------------------------- |
| `id`            | Primary key               |
| `original_link` | Original URL              |
| `code`          | Unique short code         |
| `user_id`       | Foreign key to `users.id` |
| `created_at`    | Creation timestamp        |

### users

| Column          | Description          |
| --------------- | -------------------- |
| `id`            | Primary key          |
| `name`          | Unique username      |
| `password_hash` | Hashed user password |

## Endpoints:

### Users

POST /api/users/register - create a new user
requst:

```json
{
  "username": "...",
  "password": "..."
}
```

201 success

```json
{
  "access_token": "..."
}
```

bad request 400 invalid username or password

POST /api/users/login
requst:

```json
{
  "username": "...",
  "password": "..."
}
```

response:
200 success

```json
{
  "access_token": "..."
}
```

bad request 400 invalid username or password

### Links

- GET /r/{code} - get a short link and redirect immediately without protection
  response:

302 Found
Location: `<original_url>`

404 not found:

```json
{
  "error": "not found a link"
}
```

400 bad request:

```json
{
  "error": "invalid code"
}
```

- POST /api/links - create a new short link
  request:
  only authenticated users with token

headers: bearer <token>

```json
{
  "link": "http://youtube.com/" // http or https only
}
```

response:
successfully created 201:

```json
{
  "id": 1,
  "short_link": "...",
  "created_at": "timestamp"
}
```

400 bad request

```json
{
  "error": "invalid link"
}
```

- GET /api/links - get all links user
  request:
  only authenticated users with token
  headers: bearer <token>

no json body

response:
success 200:

with links

```json
[
  {
    "id": 1,
    "short_link": "...",
    "created_at": "timestamp"
  },
  {
    "id": 2,
    "short_link": "...",
    "created_at": "timestamp"
  },
  {
    "id": 3,
    "short_link": "...",
    "created_at": "timestamp"
  }
]
```

or empty:

```json
[]
```

- GET /api/links/{id} - get a link by id

request:
only authenticated users with token
and id in params

headers: bearer <token>

no json body

response:
200:

```json
{
  "id": 1,
  "original_link": "...",
  "short_link": "...",
  "created_at": "timestamp"
}
```

400 bad request:

```json
{
  "error": "invalid id"
}
```

404 not found:

```json
{
  "error": "not found a link"
}
```

- DELETE /api/links/{id} - delete a link by id
  request:
  only authenticated users with token
  and id in params

headers: bearer <token>

no json body

response:
status code or message

404 not found:

```json
{
  "error": "not found a link"
}
```

400 bad request:

```json
{
  "error": "invalid id"
}
```

no content 204:
no json response

## Redis

For redirect GET /r/{code}:
key = link:{code}
value = original_url
TTL = 1 minute

Database (postgresql) source of truth

## Kafka

we use kafka to save how many times users visited a link for analytics
for example:

producer event link-visited:

```json
{
  "id": 42,
  "code": "x342xc"
}
```

consumer event link-visited:
just log it

```json
{
  "id": 42,
  "code": "x342xc"
}
```

"click_count": 55

## JWT authentication

claims
user_id int
and jwt.RegisteredClaims

returns access and refresh token
