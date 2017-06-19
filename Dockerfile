FROM scratch
EXPOSE 53/udp
COPY wilddns /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/wilddns"]
