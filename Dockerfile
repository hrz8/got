FROM golang:1.23.0-bookworm AS build_stage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify && go mod tidy

COPY . .
RUN go build -ldflags "-s -w" -o bin/got .

FROM gcr.io/distroless/base-debian12

COPY --from=build_stage /app/bin/got /bin/got

EXPOSE 5101 5102

ENTRYPOINT [ "/bin/got", "serve" ]
