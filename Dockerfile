FROM golang:1.25-alpine AS build
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go clean -cache -modcache -i -r
RUN CGO_ENABLED=0 GOOS=linux go build -o limitr .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest


FROM alpine:latest

WORKDIR /app
COPY --from=build /app/limitr .
COPY --from=build /go/bin/goose /usr/local/bin/goose
COPY ./sql/migrations ./sql/migrations

RUN adduser -D demo
USER demo

EXPOSE ${PORT}

CMD ["./limitr"]
