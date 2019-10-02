FROM golang:latest as builder

LABEL maintainer="walkinapp.co.uk | Jake Oliver <jake@walkinapp.co.uk>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GIN_MODE=release
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -a -installsuffix cgo -o main main.go

# Create new image and import just the binary
FROM alpine:latest

# Alpine doesn't include cert auth certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

ENV GIN_MODE=release
ENV CGO_ENABLED=0

EXPOSE 8080

CMD ["./main"]
