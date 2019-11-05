FROM scratch

EXPOSE 8080

ENV GODEBUG=netdns=go

ADD release/vela-worker /bin/

CMD ["/bin/vela-worker"]