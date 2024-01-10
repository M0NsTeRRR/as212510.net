FROM gcr.io/distroless/static:latest

COPY as212510.net /

EXPOSE 10240 10241 8080

CMD ["/as212510.net"]