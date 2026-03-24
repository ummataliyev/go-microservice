FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
RUN adduser -D -u 1000 appuser
WORKDIR /app
COPY --from=builder /bin/server /bin/server
COPY --from=builder /app/infra/.env /app/infra/.env
USER appuser
EXPOSE 8080
CMD ["/bin/server"]
