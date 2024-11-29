FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server cmd/ratelimiter/main.go

FROM scratch
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
CMD ["./server"]