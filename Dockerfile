FROM golang:1.24-bookworm AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

ENV CGO_ENABLED=1

COPY . .

RUN go build -v -o /run-app .

FROM debian:bookworm-slim

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/templates /usr/local/share/customize/templates
COPY --from=builder /usr/src/app/files /usr/local/share/customize/files

EXPOSE 8080

CMD ["run-app"]
