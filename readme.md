# URL shortener service

*You may customize service configuration in docker-compose files

## Run app with in-memory database (LRU, limit 10_000)

    docker compose -f docker-compose.yaml up --build

## Run app with postgres database

    docker compose -f docker-compose.postgres.yaml up --build

## Expected behavior

You send your url and service must return shortened id (with length 10),
then you make request to get your original url using this shortened id 

### GRPC

#### Save Request
Message

    {
        "origin_url": "https://your_url"
    }

#### Save Response
Message

    {
        "shortened_url": "0b9acd4660"
    }

#### Get Request
Message


    {
        "shortened_url": "0b9acd4660"
    }

#### Get Response
Message

    {
        "origin_url": "https://your_url"
    }

### HTTP

#### Post Request
Body

    {
        "originURL": "https://your_url"
    }

#### Post Response
Body

    {
        "shortenedURL": "0b9acd4660"
    }

#### Get Request
Url

     "http://${service_host}:{$service_port}/get/0b9acd4660"

#### Get Response
Body

    {
        "originURL": "https://your_url""
    }

## Config

```.env
ENV = DEV | PROD
PORT = 
SERVER_TYPE = 0 (GRPC) | 1 (HTTP)
DB_TYPE = 0 (POSTGRES) | 1 (INMEMORY)

# Redis config
REDIS_HOST=
REDIS_PORT=

# Postgres config
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
```

## Models

### GRPC

```protobuf
message URL {
  string origin_url = 1;
}

message ShortenedURL {
  string shortened_url = 1;
}
```

### HTTP

```go
type URL struct {
	OriginURL string `json:"originURL"`
}
```