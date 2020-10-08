FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    PORT=3000

# Move to working directory /build
# WORKDIR /build

# Copy and download dependency using go mod
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# # Copy the code into the container
# COPY . .

# RUN echo $(ls -1 ./cmd/bin-collections/)

# # Build the application
# RUN go build -o main ./cmd/bin-collections/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# # Copy binary from build to main folder
# RUN cp /build/main /build/.env.yml .

# Export necessary port
EXPOSE 3000

# Command to run when starting the container
CMD ["ls"]