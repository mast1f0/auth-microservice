FROM golang:1.25.6 AS builder

WORKDIR /auth

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/srv/main.go

FROM alpine:3.22

WORKDIR /auth

COPY --from=builder /app /auth/app

EXPOSE 8081

CMD ["./app", "-seed"]
