FROM alpine:latest

ADD webhook /bin/
RUN apk -Uuv add ca-certificates
RUN cd /bin&&ls
ENTRYPOINT /bin/webhook
