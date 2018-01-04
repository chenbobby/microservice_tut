# consignment-service/Dockerfile

### This container is used to build the binary ###

# Use official golang image, which contains correct build tools and libraries,
# as `builder`, which is the container name that we can reference later
FROM golang:latest as builder

# Set WORKDIR to our current service in the gopath
WORKDIR /go/src/github.com/chenbobby/microservice_tut/consignment-service

# Copy current code into WORKDIR
COPY . .

# Install `godep` inside container to deal with dependencies
RUN go get -u github.com/golang/dep/cmd/dep

# Create a `dep` project, and run `dep ensure` to pull in dependencies
RUN dep init && dep ensure

# Build the binary, with extra flags so the binary will run in Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


### This container is used to run the binary ###

# This 2nd FROM statement tells Docker to start a new build process with this image
FROM alpine:latest

# Security related package
RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

# Instead of copying the binary from our local host machine
# ADD consignment-service /app/consignment-service
# We pull the binary from the container named `builder` (see above)
COPY --from=builder /go/src/github.com/chenbobby/microservice_tut/consignment-service/consignment-service .

CMD ["./consignment-service"]