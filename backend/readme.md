# Labbaika Payslip

## Run Service & Worker
### run user service
`go run cmd/main.go payslip-svc`

### run worker
`go run cmd/main.go ??`

## run all service & worker with docker-compose
`docker-compose up`

## Build Application
`make build` or `go build -o labbaika ./cmd`

## Folder structure
```
/cmd
  /app                      # Main application code (eg. user api & worker)
  /config                   # Configuration files
/docs                       # Docs related to this service (API, etc)
/internal
  /app                      # business logic
    /domain                 # Domain entities
    /repository             # Interfaces defining repository contracts
    /usecase                # Use cases or interactors
/localdevscripts            # Contain localstack script to init SNS topic & SQS queue
/pkg
  /database                 # Database related code (repositories implementation, migrations, etc.)
    /migrations             # SQL script for migrations
  /errors                   # Errors library & handling
  /external                 # External service clients
  /http                     # HTTP delivery mechanism
    /handler                # HTTP request handlers
    /middleware             # Middleware functions
/mocks                      # Generated mock implementations of interfaces (used for testing)
```

