# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang AS build

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
FROM alpine AS release
RUN adduser --system --uid 1234 figment
WORKDIR /app/kava
COPY --from=build /go/src/github.com/figment-networks/worker-kava/worker /app/kava/worker
RUN chmod a+x ./worker
RUN chown -R figment /app/
USER 1234
CMD ["./worker"]
