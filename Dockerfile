FROM alpine:latest
EXPOSE 3000

ADD drone-webhook /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/drone-webhook
