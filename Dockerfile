#FROM golang:1.24
#ENV GOPATH=/
#
#COPY ./ ./
#
#RUN go mod download
#RUN go build -o auth-app ./cmd/main.go
#
#ENTRYPOINT ["./auth-app"]

FROM golang:1.24 AS builder
ENV GOPATH=/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth ./cmd/main.go


FROM alpine:latest
WORKDIR /
COPY --from=builder /auth /auth
COPY .env .
ENTRYPOINT ["./auth"]