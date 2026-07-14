FROM golang:1.26.4@sha256:f96cc555eb8db430159a3aa6797cd5bae561945b7b0fe7d0e284c63a3b291609 AS builder

ARG VERSION=development
ARG SOURCE_DATE_EPOCH=0

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -a -o as212510.net -ldflags "-w -X main.version=$VERSION -X main.buildTime=$SOURCE_DATE_EPOCH -extldflags '-static'" cmd/as212510/main.go

FROM gcr.io/distroless/static:nonroot@sha256:f7f8f729987ad0fdf6b05eeeae94b26e6a0f613bdf46feea7fc40f7bd72953e6

COPY --from=builder /go/src/app/as212510.net /bin/as212510.net

USER 65532

EXPOSE 8080 10240 10241

CMD ["/bin/as212510.net"]
