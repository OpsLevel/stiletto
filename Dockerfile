FROM golang:1.20

COPY ./src /src
WORKDIR /src
RUN go build -o /bin/stiletto . && chmod +x /bin/stiletto
ENTRYPOINT ["/entrypoint.sh"]
COPY entrypoint.sh /entrypoint.sh