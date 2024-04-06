FROM gcr.io/distroless/static:nonroot@sha256:f41b84cda410b05cc690c2e33d1973a31c6165a2721e2b5343aab50fecb63441

COPY as212510.net /

EXPOSE 10240 10241 8080

CMD ["/as212510.net"]