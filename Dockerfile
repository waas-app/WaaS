FROM node:18-alpine as website

WORKDIR /code
COPY ./website/package.json ./
COPY ./website/yarn.lock ./
RUN yarn
COPY ./website/ ./
RUN yarn build

FROM golang:alpine as builder
LABEL org.opencontainers.image.source https://github.com/waas-app/WaaS

RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./
COPY . .
RUN apk add curl
RUN curl -sSL https://github.com/maxcnunes/waitforit/releases/download/v2.4.1/waitforit-linux_amd64 -o wait-for-it

RUN apk add protobuf
RUN apk add protobuf-dev
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN export PATH="$PATH:$(go env GOPATH)/bin"
# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container
RUN protoc --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/devices.proto
RUN protoc --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/server.proto
# Build the Go app
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=1
ENV GO111MODULE=on
RUN go build -o waas main.go
RUN go build -o worker cmd/worker/main.go
RUN ls -aril

# Start a new stage from scratch
FROM alpine:latest as waas
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/waas .
# COPY --from=builder /app/waas.yml .
COPY --from=builder /app/wait-for-it /usr/local/bin/
COPY --from=website /code/build /website/build
RUN ls -aril
# RUN cat waas.yml
RUN chmod +x /usr/local/bin/wait-for-it

# Expose port 8080 to the outside world
EXPOSE 8000

RUN chmod +x waas
RUN apk add iptables
RUN apk add wireguard-tools
RUN apk add curl

#Command to run the executable
# CMD ["./waas", "serve"]
CMD wait-for-it -host=postgresql -port=5432 -timeout=60 -- ./waas serve

FROM alpine:latest as worker
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/worker .
COPY --from=builder /app/wait-for-it /usr/local/bin/
RUN ls -aril
# RUN cat waas.yml
RUN chmod +x /usr/local/bin/wait-for-it

# CMD [ "./worker" ]
CMD wait-for-it -host=redis -port=6379 -timeout=60 -- ./worker