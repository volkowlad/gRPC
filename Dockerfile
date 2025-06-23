FROM golang:1.24 AS builder
ENV GOPATH=/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go


FROM alpine:latest
WORKDIR /
COPY --from=builder /auth-service /auth-service
COPY .env .
ENTRYPOINT ["./auth-service"]