FROM golang:alpine AS builder
MAINTAINER cylon
WORKDIR /uranus
COPY ./ /uranus
ENV GOPROXY https://goproxy.cn,direct
RUN \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add upx bash make && \
    make build && \
    upx -1 _output/firewalld-gateway && \
    chmod +x _output/firewalld-gateway

FROM alpine AS runner
WORKDIR /uranus
COPY --from=builder /uranus/_output/firewalld-gateway /usr/sbin/
COPY --from=builder /uranus/firewalld-gateway.toml .
RUN firewalld-gateway --sql-driver=sqlite --migration
VOLUME ["/uranus" ]
CMD ["firewalld-gateway", "--sql-driver=sqlite", "--config", "firewalld-gateway.toml", "-v", "10"]
EXPOSE 2952/tcp