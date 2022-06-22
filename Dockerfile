FROM golang:1.18.3

RUN mkdir /logsrv-config

WORKDIR /logsrv

# RUN apt-get update
# RUN apt-get install -y htop mc tilde nano

COPY ./ .
RUN make build

EXPOSE 8080

ENTRYPOINT ["./logsrv", "-config-path", "/logsrv-config/config.toml"] 