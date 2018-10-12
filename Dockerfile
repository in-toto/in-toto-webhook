FROM golang:1.10 as builder

RUN mkdir -p /go/src/github.com/santiagotorres/in-toto-webhook/

WORKDIR /go/src/github.com/santiagotorres/in-toto-webhook

COPY . .

RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*") && \
  GIT_COMMIT=$(git rev-list -1 HEAD) && \
  CGO_ENABLED=0 GOOS=linux go build \
  -a -installsuffix cgo -o in-toto ./cmd/in-toto

FROM intoto/base:latest

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    ca-certificates

WORKDIR /home/app

COPY --from=builder /go/src/github.com/santiagotorres/in-toto-webhook/in-toto .
COPY certs certs

RUN chown -R app:app ./

USER app

CMD ["./in-toto"]
