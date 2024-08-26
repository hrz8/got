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
