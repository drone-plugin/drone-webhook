FROM alpine:latest
EXPOSE 3000

ADD webhook /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/webhook
