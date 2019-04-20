FROM alpine:latest

# 用于标示创建者邮箱, 协同项目不需要
LABEL maintainer="zhensheng.five@gmail.com"

# 修正时区为东8区
RUN apk add --no-cache --virtual .build-deps \
        tzdata \
        && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
        && echo "Asia/Shanghai" > /etc/timezone \
        && apk del .build-deps

ENV TZ "Asia/Shanghai"

# 安装ca证书
RUN apk add --update --no-cache \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY identifier /usr/local/bin/identifier
COPY docker-entrypoint /usr/local/bin/

WORKDIR /usr/local/var/identifier

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
