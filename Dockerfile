# Default to Go 1.21
ARG GO_VERSION=1.23

# Start from golang v1.18 base image
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# 'ca-certificates' is a package containing a set of SSL/TLS certificates that are trusted by the system.

RUN apk add --no-cache ca-certificates


# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /meal-management-backend

# Import the code from the context.
COPY ./ ./


# Build the Go app
RUN GOFLAGS=-mod=mod GOOS=linux  go build -ldflags="-w -s" -a -o  /app .

######## Start a new stage from scratch #######
# Final stage: the running container.
FROM alpine:3.16  AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable from the first stage.

WORKDIR /meal-management-backend

COPY --from=builder /app /app


RUN apk --no-cache add tzdata

# Create the directory for photo storage
RUN mkdir -p /tmp/photos

# Ensure proper permissions for the photos directory
RUN chown nobody:nobody /tmp/photos

# Use the unprivileged user for security
USER nobody:nobody

# Perform any further action as an unprivileged user.
#USER nobody:nobody

EXPOSE 64000

# Run the compiled binary.
ENTRYPOINT ["/app"]