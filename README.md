# SumUp Bank API

A Go-based API server for managing bank accounts.

## Prerequisites

- Go 1.25.7 or later
- Docker or Podman (for running the database locally)

## Testing
First make sure the db is running & migrated by running:
```bash
make up
```
Then run the tests:
```bash
go test -v ./...
```

## Running in Development

### 1. Start the Database

Start PostgreSQL and run migrations using Docker/Podman:

```bash
make up
```

This will:
- Start a PostgreSQL 16 container on port 5432
- Run database migrations automatically

To stop the database:

```bash
make down
```

### 2. Run the API Server

```bash
go run . serve
```

### 3. Open Documentation in Browser (Optional)

Once the API server is running, you can view the API documentation:

```bash
open http://localhost:8080/api/docs
```
**Note**: you can try the API using the UI by clicking "Try it out" button :) 

## API Documentation

The API is defined using OpenAPI 3.0 specification. The specification file is located at `api/openapi.yaml`.

### Generate/Update API Code

We use [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) to generate Go server code from the OpenAPI specification.

Install the code generator tool:

```bash
go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

After modifying the OpenAPI specification (`api/openapi.yaml`), regenerate the Go server code:

```bash
go generate ./api/...
```
