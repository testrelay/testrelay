FROM golang:alpine as builder


WORKDIR /app
COPY . .

RUN go build -o main cmd/server/main.go

FROM alpine:latest

# install certs and git as we need these in the base images
# to make sure that we can clone the test repos
RUN apk --no-cache add ca-certificates && \
    apk update && apk upgrade && \
    apk add --no-cache bash git openssh

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]