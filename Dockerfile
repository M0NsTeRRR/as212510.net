FROM gcr.io/distroless/static:latest

COPY as212510.net /app/as212010.net

EXPOSE 10240 8080

ENTRYPOINT ["/app/as212510.net"]