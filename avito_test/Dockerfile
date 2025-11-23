FROM golang:1.24.1-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o myapp ./cmd/avito_test/main.go


FROM alpine:3.22.2

WORKDIR /app


COPY --from=builder /app/myapp .

CMD ["./myapp"]

EXPOSE 8080