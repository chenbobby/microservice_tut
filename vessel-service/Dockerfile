# vessel-service/Dockerfile

### This container is used to build the binary ###
FROM golang:latest as builder

WORKDIR /go/src/github.com/chenbobby/microservice_tut/vessel-service

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a installsuffix cgo .


### This container is used to run the binary ###
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/github.com/chenbobby/microservice_tut/vessel-service/vessel-service .

CMD ["./vessel-service"]