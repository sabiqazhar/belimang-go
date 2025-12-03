FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary dari builder
COPY --from=builder /app/main .

# Expose port sesuai service
EXPOSE 8080

# Jalankan binary
CMD ["./main"]