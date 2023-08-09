FROM golang:1.20

COPY ./src /src
WORKDIR /src
ENTRYPOINT ["/entrypoint.sh"]
COPY entrypoint.sh /entrypoint.sh