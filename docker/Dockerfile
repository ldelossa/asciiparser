# Multi-stage build. See: https://docs.docker.com/develop/develop-images/multistage-build/
# Builds the go binary in a container with the runtime then copies to a stripped down container
FROM golang:latest
WORKDIR /go/src/github.com/ldelossa/asciiparser
COPY . .
RUN ls ./vendor #confirm vendor folder dependencies are copied 
# build confd
RUN cd ./cmd/asciiparser && CGO_ENABLED=0 GOOS=linux go build -a -o asciiparser .

# second stage. copy built binary to alpine container and set entry point to the compiled binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/ldelossa/asciiparser/cmd/asciiparser/asciiparser .
CMD ["./asciiparser"]
