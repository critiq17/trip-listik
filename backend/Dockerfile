FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/migrate ./cmd/migrate

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /bin/api /app/api
COPY --from=builder /bin/migrate /app/migrate
EXPOSE 8080
ENTRYPOINT ["/app/api"]
