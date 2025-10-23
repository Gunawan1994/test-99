FROM golang:1.24 AS builder
RUN apt-get update -y
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app
COPY . .

RUN [ -f go.mod ] && rm go.mod || echo "go.mod not found, skip"
RUN [ -f go.sum ] && rm go.sum || echo "go.sum not found, skip"

RUN go mod init listing
RUN go mod tidy

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /go/bin/app

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
COPY --from=builder /go/bin/app .
ENTRYPOINT ["./app"]
