package main

const envTemplate = `# Basic env configuration

# Service Handlers
## Server
USE_REST=true
USE_GRPC=true
USE_GRAPHQL=true
## Worker
USE_KAFKA_CONSUMER=true
USE_CRON_SCHEDULER=true
USE_REDIS_SUBSCRIBER=true

REST_HTTP_PORT=8000
GRAPHQL_HTTP_PORT=8001
GRPC_PORT=8002

BASIC_AUTH_USERNAME=user
BASIC_AUTH_PASS=pass

MONGODB_HOST_WRITE=mongodb://localhost:27017
MONGODB_HOST_READ=mongodb://localhost:27017
MONGODB_DATABASE_NAME={{.ServiceName}}

SQL_DRIVER_NAME=[string]
SQL_DB_READ_HOST=[string]
SQL_DB_READ_USER=[string]
SQL_DB_READ_PASSWORD=[string]
SQL_DB_WRITE_HOST=[string]
SQL_DB_WRITE_USER=[string]
SQL_DB_WRITE_PASSWORD=[string]
SQL_DATABASE_NAME=[string]

REDIS_READ_HOST=localhost
REDIS_READ_PORT=6379
REDIS_READ_AUTH=
REDIS_WRITE_HOST=localhost
REDIS_WRITE_PORT=6379
REDIS_WRITE_AUTH=

KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID={{.ServiceName}}
KAFKA_CONSUMER_GROUP={{.ServiceName}}

JAEGER_TRACING_HOST=[string]
GRAPHQL_SCHEMA_DIR="api/{{.ServiceName}}/graphql"
JSON_SCHEMA_DIR="api/{{.ServiceName}}/jsonschema"


# Additional env

`
