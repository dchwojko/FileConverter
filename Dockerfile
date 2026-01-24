# Build stage - use Ubuntu base
FROM ubuntu:22.04 AS builder

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Install Go
RUN apt-get update && apt-get install -y \
    wget \
    tar \
    && wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz \
    && rm go1.21.6.linux-amd64.tar.gz \
    && rm -rf /var/lib/apt/lists/*

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage - use minimal Ubuntu
FROM ubuntu:22.04

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]