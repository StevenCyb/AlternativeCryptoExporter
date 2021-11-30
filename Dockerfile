# Build executable binary
FROM golang:alpine3.13 AS builder

ENV USER=authg
ENV UID=1000

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "$(pwd)" \
  --no-create-home \
  --uid "1000" \
  "ace"

RUN apk --update add ca-certificates

WORKDIR /build

COPY go/go.mod .
COPY go/go.sum .

RUN go mod download
RUN go mod verify

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go/ .

RUN go build -o ace .

WORKDIR /dist

RUN cp /build/ace .

# Build image
FROM scratch

LABEL Name=ace \
      Release=https://github.com/StevenCyb/AlternativeCryptoExporter \
      Url=https://github.com/StevenCyb/AlternativeCryptoExporter \
      Help=https://github.com/StevenCyb/AlternativeCryptoExporter/issues

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /dist/ace /

USER ace:ace
ENTRYPOINT ["/ace"]