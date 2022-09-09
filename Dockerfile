FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN echo "deb http://mirrors.aliyun.com/debian stretch main contrib non-free" > /etc/apt/sources.list \
        && echo "deb-src http://mirrors.aliyun.com/debian stretch main contrib non-free" >> /etc/apt/sources.list \
        && echo "deb http://mirrors.aliyun.com/debian stretch-updates main contrib non-free" >> /etc/apt/sources.list \
        && echo "deb-src http://mirrors.aliyun.com/debian stretch-updates main contrib non-free" >> /etc/apt/sources.list \
        && echo "deb http://mirrors.aliyun.com/debian-security stretch/updates main contrib non-free" >> /etc/apt/sources.list \
        && echo "deb-src http://mirrors.aliyun.com/debian-security stretch/updates main contrib non-free" >> /etc/apt/sources.list \
        && apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app/bin
COPY --from=builder /src/conf /app/conf
COPY --from=builder /src/static /app/static
COPY --from=builder /src/storage /app/storage

WORKDIR /app

EXPOSE 8888
VOLUME /app/conf
VOLUME /app/storage

CMD ["./bin/app", "--conf", "/app/conf/config.yaml"]