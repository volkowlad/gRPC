FROM golang:1.24
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o tasks-app ./cmd/main.go

ENTRYPOINT ["./tasks-app"]

#FROM golang:1.24 AS builder
#ENV GOPATH=/
#COPY . .
#RUN go mod download
#RUN --mount=type=secret,id=dotenv,target=.env
#RUN CGO_ENABLED=0 GOOS=linux go build -o /auth ./cmd/main.go
#RUN rm -f .env
#
#
#FROM alpine:latest
#WORKDIR /
#COPY --from=builder /auth /auth
#ENTRYPOINT ["/auth"]