FROM golang:1.23.3 AS builder
RUN apt-get update -y

# Environment variables
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOTOOLCHAIN=auto

WORKDIR /go/src/app
COPY . .

# Optional: hapus go.mod/sum kalau mau auto init
RUN [ -f go.mod ] && rm go.mod || echo "go.mod not found, skip"
RUN [ -f go.sum ] && rm go.sum || echo "go.sum not found, skip"

RUN go mod init user
RUN go mod tidy

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /go/bin/app

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
COPY --from=builder /go/bin/app .
ENTRYPOINT ["./app"]
