FROM golang:1.25 AS builder
WORKDIR /blog-src
COPY main.go .
RUN go mod init blog
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o blog-server main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev ca-certificates procps
WORKDIR /blog
COPY --from=builder /blog-src/blog-server .
COPY templates ./templates
COPY main_test.go ./
#CMD ["sleep", "infinity"]
ENTRYPOINT ["./blog-server"]
