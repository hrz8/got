# got

## setup

```bash
cp .env.example .env # change with yours
go mod tidy
export $(cat .env | xargs) && go run . migrate up
```

## run

```bash
export $(cat .env | xargs) && go run . serve
```

## build

```bash
go build -ldflags "-s -w" -o bin/got . 
```

## run binary

```bash
./bin/got serve
# with .env
export $(cat .env | xargs) && ./bin/got serve
```

## dockerize

### build
```bash
docker build -t got:1.0.0 .
```

### run
```bash
docker run -d \
-e DATABASE_URL='postgresql://postgres:pass@localhost:5432/got' \
-e DATABASE_URL_READER='postgresql://postgres:pass@localhost:5432/go' \
-e DATABASE_NAME='got' \
-e LOG_LEVEL='debug' \
-p 5001:5101 \
-p 5002:5102 \
got:1.0.0
```

## exe

### http
```bash
curl http://localhost:5101/healthz 
curl -X POST http://localhost:5101/v1/greeter/hello -H 'Content-Type: application/json' -d '{"name": "john"}'
# chi router
curl http://localhost:5101/v1/users -H 'X-Rest-Request: true'
curl http://localhost:5101/v1/users  # will return 404 if no X-Rest-Request
# extras gwmux
curl curl http://localhost:5101/users/123
```

### grpc
```bash
grpcurl -plaintext localhost:5102 grpc.health.v1.Health/Check
grpcurl -plaintext -d '{"name": "John"}' localhost:5102 service.v1.GreeterService/SayHello
```
