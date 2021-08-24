FROM golang:1.16 as builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o fileapi cmd/http/http.go

FROM gcr.io/distroless/base
COPY  --from=builder /app/fileapi .

CMD ["./fileapi"]
