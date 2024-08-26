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
