FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

ENV CGO_ENABLED=1

COPY . .

RUN go build -v -o /app/customize .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/customize /app/
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/files /app/files

EXPOSE 8080

CMD ["./customize"]