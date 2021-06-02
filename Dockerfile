# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.14 AS build

WORKDIR /go/src/github.com/figment-networks/worker-kava/

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY .git .git
COPY ./Makefile ./Makefile
COPY ./api ./api
COPY ./client ./client
COPY ./cmd/common ./cmd/common
COPY ./cmd/worker-kava ./cmd/worker-kava


ENV GOARCH=amd64
ENV GOOS=linux

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10 AS release

WORKDIR /app/kava
COPY --from=build /go/src/github.com/figment-networks/worker-kava/worker /app/kava/worker
RUN chmod a+x ./worker
CMD ["./worker"]
