FROM golang:1.26.0@sha256:c83e68f3ebb6943a2904fa66348867d108119890a2c6a2e6f07b38d0eb6c25c5 AS builder

ARG VERSION=development
ARG SOURCE_DATE_EPOCH=0

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -a -o as212510.net -ldflags '-w -X main.version=$VERSION -X main.buildTime=$SOURCE_DATE_EPOCH -extldflags "-static"'

FROM gcr.io/distroless/static:nonroot@sha256:cba10d7abd3e203428e86f5b2d7fd5eb7d8987c387864ae4996cf97191b33764

COPY --from=builder /go/src/app/as212510.net /bin/as212510.net

USER 65532

EXPOSE 8080 10240 10241

CMD ["/bin/as212510.net"]
