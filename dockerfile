FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .


RUN go build -o main .

# ====== Stage 2: Runtime ======

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /app/db/schema /app/db/schema
COPY --from=builder /app/db/queries /app/db/queries
COPY --from=builder /app/main .
COPY --from=builder /app/web/ /app/web

ENV GOOSE_DRIVER=postgres
ENV GOOSE_MIGRATION_DIR=./db/schema

EXPOSE ${PORT}

CMD ["./main"]
