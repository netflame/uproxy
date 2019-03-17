# keep source image same as scrapd.dockerfile (Dockerfile)
FROM alpine:3.8

LABEL maintainer=ilpan:<pna.dev@outlook.com>

WORKDIR /uproxy

# make sure have this souce before copy it
COPY ./uproxy_linux_amd64 ./uproxy
COPY ./config.toml .
# COPY ./sites.toml .
COPY ./docker-entrypoint.sh .

RUN sed -i "s#localhost:#redis:#" ./config.toml

EXPOSE 6001

ENTRYPOINT [ "./docker-entrypoint.sh", "./uproxy"]