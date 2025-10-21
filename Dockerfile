ARG GO_VERSION=1.25-bookworm

# Build Goose
FROM golang:${GO_VERSION} AS build_goose
ARG GOOSE_VERSION=v3.25.0
RUN apt-get update && apt-get install git -y
WORKDIR /goose_src
RUN git clone --depth 1 --branch ${GOOSE_VERSION} https://github.com/pressly/goose.git .
RUN go build -tags="no_mysql no_sqlite3 no_ydb no_clickhouse no_mssql no_vertica" -o /go/bin/goose ./cmd/goose

# Build Limitr
FROM golang:${GO_VERSION} AS build_limitr
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w -buildid=" -o limitr .

# Run migrations
FROM golang:${GO_VERSION} AS migrate_runner
WORKDIR /migrate
COPY --from=build_goose /go/bin/goose /usr/local/bin/goose
COPY ./sql/migrations ./sql/migrations

# Start Limitr
FROM gcr.io/distroless/base-nossl AS final
WORKDIR /app
COPY --from=build_limitr /app/limitr .

USER nobody

ARG PORT=8080
ENV PORT=${PORT}
EXPOSE ${PORT}

ENTRYPOINT ["/app/limitr"]
