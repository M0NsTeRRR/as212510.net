FROM golang:1.25.5@sha256:36b4f45d2874905b9e8573b783292629bcb346d0a70d8d7150b6df545234818f AS builder

ARG VERSION=development
ARG SOURCE_DATE_EPOCH=0

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -a -o as212510.net -ldflags '-w -X main.version=$VERSION -X main.buildTime=$SOURCE_DATE_EPOCH -extldflags "-static"'

FROM gcr.io/distroless/static:nonroot@sha256:2b7c93f6d6648c11f0e80a48558c8f77885eb0445213b8e69a6a0d7c89fc6ae4

COPY --from=builder /go/src/app/as212510.net /bin/as212510.net

USER 65532

EXPOSE 8080 10240 10241

CMD ["/bin/as212510.net"]
