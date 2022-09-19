FROM alpine:latest

ADD webhook /bin/
RUN apk -Uuv add ca-certificates
RUN ls
ENTRYPOINT /bin/webhook
