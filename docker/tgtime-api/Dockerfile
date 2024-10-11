FROM busybox:1.35.0-uclibc AS busybox

FROM mockserver/mockserver:5.15.0

COPY --from=busybox /bin /bin

USER root

RUN wget -O curl.tar.xz https://github.com/stunnel/static-curl/releases/download/8.2.1/curl-static-amd64-8.2.1.tar.xz &&  \
    tar -xf curl.tar.xz --directory /usr/bin

CMD ["bin/sh"]