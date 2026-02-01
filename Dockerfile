# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Run stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]
