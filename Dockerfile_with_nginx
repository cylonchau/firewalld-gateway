FROM golang:alpine AS builder
MAINTAINER cylon
WORKDIR /uranus
COPY ./ /uranus
ENV GOPROXY https://goproxy.cn,direct
RUN \
    #sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add upx bash make && \
    make build && \
    upx -1 _output/firewalld-gateway && \
    chmod +x _output/firewalld-gateway

FROM nginx:1.20 AS runner
WORKDIR /uranus
ARG S6_OVERLAY_VERSION=3.1.5.0
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64.tar.xz /tmp

RUN apt update && apt install xz-utils procps iproute2 -y && \
    tar -Jxpf /tmp/s6-overlay-x86_64.tar.xz -C / && \
    tar -Jxpf /tmp/s6-overlay-noarch.tar.xz -C / && \
    rm -f /tmp/s6-overlay-x86_64.tar.xz && \
    rm -f /tmp/s6-overlay-noarch.tar.xz
ENTRYPOINT ["/init"]
RUN mkdir /etc/services.d/
COPY --from=builder /uranus/_output/firewalld-gateway ./bin/
COPY --from=builder /uranus/firewalld-gateway.toml .
COPY --from=builder /uranus/dist /var/run/dist/
COPY --from=builder /uranus/uranus.nginx-2953.conf /etc/nginx/conf.d/fw.conf
COPY --from=builder /uranus/s6/ /etc/s6-overlay/s6-rc.d/
COPY --from=builder /uranus/s6/ /etc/services.d/
ENV PATH "$PATH:/uranus/bin"
RUN  firewalld-gateway --sql-driver=sqlite --migration && \
     rm -f /etc/nginx/conf.d/default.conf && \
     echo "longrun" > /etc/s6-overlay/s6-rc.d/nginx/type && \
     echo "longrun" > /etc/s6-overlay/s6-rc.d/uranus/type && \
     mkdir -pv /etc/s6-overlay/s6-rc.d/uranus/contents.d && \
     mkdir -pv /etc/s6-overlay/s6-rc.d/nginx/contents.d

#CMD [ " /command/s6-svscan", "/etc/services.d" ]
VOLUME ["/uranus" ]
EXPOSE 2953/tcp