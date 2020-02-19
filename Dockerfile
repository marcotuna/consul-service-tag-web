# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.13-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the command inside the container.
RUN make

# Use a Docker multi-stage build to create a lean production image.
FROM gcr.io/distroless/base-debian10

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/consul-service-tag-web /consul-service-tag-web

# Copy templates, static and such
COPY template template/

# Run the web service on container startup.
CMD ["/consul-service-tag-web"]