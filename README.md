# Todo Project for graph

## Table of Contents

- [Setup](#setup)
- [Swagger](#swagger)
- [Tests](#tests)
- [Prometheus](#prometheus)
- [LoadTest](#loadtest)



## Setup

To set up the project, first copy `example.env` and create an `.env` file in the root directory. Do the same for `example.config.yaml` by copying it to `config.yaml`.

> **Note:** You can change the `port`, `password`, `dbname`, and other settings, but make sure to update them in both `config.yaml` and `.env`.

```bash
cp example.env .env
cp example.config.yaml config.yaml
```
> after this run docker compose

```bash
docker compose up -d
```

TODO

---

## Swagger

After building and starting the application with Docker Compose, you can access the Swagger API documentation at:

```text
http://localhost:<port>/swagger/index.html
http://localhost:8080/swagger/index.html
```

```bash
create new task request

curl -X 'POST' \
  'http://localhost:8080/api/v1/todo/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "assignee": "amir",
  "description": "this is simple task",
  "priority": 1,
  "title": "test task"
}'

response

{
  "status": "success",
  "code": 201,
  "message": "Todo created successfully",
  "data": {
    "uuid": "f3253ae8-0d19-4d6c-9321-3ccfc644280b"
  },
  "errors": null
}

---------------------

request delete task


curl -X 'DELETE' \
  'http://localhost:8080/api/v1/todo/delete/f3253ae8-0d19-4d6c-9321-3ccfc644280b' \
  -H 'accept: application/json'


response

{
  "status": "success",
  "code": 200,
  "message": "Todo deleted successfully",
  "data": null,
  "errors": null
}

---------------------

request to get list of tasks

curl -X 'GET' \
  'http://localhost:8080/api/v1/todo/list?page=1&status=0&priority=1' \
  -H 'accept: application/json'

response

{
  "status": "success",
  "code": 200,
  "message": "Todos retrieved successfully",
  "data": {
    "todos": [
      {
        "uuid": "76766e31-4664-4f63-bed1-1566ce05601e",
        "assignee": "amir",
        "title": "test task",
        "description": "this is simple task",
        "priority": 1,
        "status": 0,
        "todo_start_at": "2026-09-11T09:48:41.525434Z",
        "do_start_at": null,
        "done_start_at": null,
        "created_at": "2026-09-11T09:48:41.525434Z",
        "updated_at": "2026-09-11T09:48:41.525434Z"
      }
    ],
    "total": 1
  },
  "errors": null
}

---------------------

request to update task

curl -X 'PUT' \
  'http://localhost:8080/api/v1/todo/update/76766e31-4664-4f63-bed1-1566ce05601e' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "assignee": "mohamad",
  "description": "update description",
  "priority": 1,
  "status": 1,
  "title": "update title"
}'

response

{
  "status": "success",
  "code": 200,
  "message": "Update successfully",
  "data": {
    "uuid": "76766e31-4664-4f63-bed1-1566ce05601e",
    "assignee": "mohamad",
    "title": "update title",
    "description": "update description",
    "priority": 1,
    "status": 1,
    "todo_start_at": "2026-09-11T09:48:41.525434Z",
    "do_start_at": "2026-09-11T09:50:05.881470083Z",
    "done_start_at": null,
    "created_at": "2026-09-11T09:48:41.525434Z",
    "updated_at": "2026-09-11T09:48:41.525434Z"
  },
  "errors": null
}

---------------------


request to get info task

curl -X 'GET' \
  'http://localhost:8080/api/v1/todo/76766e31-4664-4f63-bed1-1566ce05601e' \
  -H 'accept: application/json'

response

{
  "status": "success",
  "code": 200,
  "message": "Todo retrieved successfully",
  "data": {
    "todo": {
      "uuid": "76766e31-4664-4f63-bed1-1566ce05601e",
      "assignee": "mohamad",
      "title": "update title",
      "description": "update description",
      "priority": 1,
      "status": 1,
      "todo_start_at": "2026-09-11T09:48:41.525434Z",
      "do_start_at": "2026-09-11T09:50:05.88147Z",
      "done_start_at": null,
      "created_at": "2026-09-11T09:48:41.525434Z",
      "updated_at": "2026-09-11T09:50:05.883022Z"
    }
  },
  "errors": null
}


-----------------

request create validation error


curl -X 'POST' \
  'http://localhost:8080/api/v1/todo/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{

}'

{
  "status": "error",
  "code": 400,
  "message": "Validation failed",
  "data": null,
  "errors": [
    {
      "title": "title is required"
    },
    {
      "assignee": "assignee is required"
    },
    {
      "description": "description is required"
    },
    {
      "priority": "priority is required"
    }
  ]
}

curl -X 'POST' \
  'http://localhost:8080/api/v1/todo/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "assignee": "amir",
  "description": "this is simple task",
  "priority": 5,
  "title": "test task"
}'

{
  "status": "error",
  "code": 400,
  "message": "Validation failed",
  "data": null,
  "errors": [
    {
      "priority": "priority must be between 1 and 3"
    }
  ]
}

curl -X 'POST' \
  'http://localhost:8080/api/v1/todo/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "assignee": "amir",
  "description": "this is simple task",
  "priority": 2,
  "title": "test task"
}'

{
  "status": "error",
  "code": 409,
  "message": "todo with this title already exists",
  "data": null,
  "errors": null
}

------------------------

validation error for delete

curl -X 'DELETE' \
  'http://localhost:8080/api/v1/todo/delete/f3253ae8-0d19-4d6c-9321-3ccfc6442801' \
  -H 'accept: application/json'


{
  "status": "error",
  "code": 404,
  "message": "todo not found",
  "data": null,
  "errors": null
}

```

---

## Tests

The project test coverage is focused on the service layer because all business logic is implemented there.  
As a result, testing the service layer provides sufficient coverage of the application's core logic, and other layers do not require additional unit tests.

The current test coverage is **88%**.

![Test Coverage](/assets/test-coverage.png)

---

## Prometheus

You can view the Prometheus metrics at:

```text
http://localhost:<port>/metrics
You can also access the Prometheus dashboard at:
http://localhost:9090
each 15 second gather data from metrics
```
From the Prometheus dashboard, you can search and query the available metrics.

![Prometheus](/assets/prometheus.png)

---

## LoadTest

I use the `hey` tool to perform a load test with **10,000 requests** and a **concurrency of 5**.

```bash
hey -n 10000 -c 5 http://localhost:8080/api/v1/todo/list
```

At the same time, I run `pprof` to collect a 30-second CPU profile:
```bash
go tool pprof "http://localhost:8080/debug/pprof/profile?seconds=30"
```
While pprof is collecting the profile, I run the hey load test:
```bash
hey -n 10000 -c 5 http://localhost:8080/api/v1/todo/list
```
# this is analyze heap
![heap](/assets/heap.png)

# this is analyze allocations
![heap](/assets/allocs.png)


---