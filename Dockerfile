FROM golang:1.25.6
WORKDIR /auth
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app cmd/srv/main.go
EXPOSE 8081
CMD ["./app"]