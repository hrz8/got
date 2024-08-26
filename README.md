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
curl -X POST http://localhost:5101/v1/greeter/hello -d '{"name": "john"}'
```

### grpc
```bash
grpcurl -plaintext localhost:5102 service.v1.Health/Check
grpcurl -plaintext -d '{"name": "john"}' localhost:5102 service.v1.Greeter/SayHello
```
