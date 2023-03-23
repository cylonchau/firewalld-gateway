FROM golang:alpine AS builder
MAINTAINER cylon
WORKDIR /firewall
COPY ./ /firewall
ENV GOPROXY https://goproxy.cn,direct
RUN \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add upx bash make && \
    make build && \
    upx -1 _output/firewalld-gateway && \
    chmod +x _output/firewalld-gateway

FROM alpine AS runner
WORKDIR /go/firewalld
COPY --from=builder /firewalld/_output/firewalld-gateway ./bin/
COPY --from=builder /firewalld/firewalld-gateway.conf .
CMD  ["firewalld-gateway", "-v", "5", "--config", "./firewalld-gateway.conf"]
ENV PATH "$PATH:/go/firewalld/bin"
VOLUME ["/firewall"]