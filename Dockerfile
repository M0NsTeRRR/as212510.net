FROM golang:1.25.6@sha256:06d1251c59a75761ce4ebc8b299030576233d7437c886a68b43464bad62d4bb1 AS builder

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
