# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.17 AS builder

# Copy local code to the container image.
WORKDIR /go/src/app
COPY . .

# Fetch dependencies.
# Using go get.
RUN go mod download

# Build the command inside the container.
# Disable CGO to make the binary static for better compatibility.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o app

# Use a minimal image to run the compiled Go binary
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/app/app /app

# Run the web service on container startup.
CMD ["/app"]