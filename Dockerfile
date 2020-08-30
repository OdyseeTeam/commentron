## Mooney API's for LBRY Inc.
FROM ubuntu:18.04
LABEL MAINTAINER="beamer"

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get -yq install apt-utils tzdata wait-for-it ca-certificates && \
    apt-get autoclean -y && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /usr/bin
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
COPY ./bin/mooney mooney
RUN chmod +x ./mooney

EXPOSE 6060
STOPSIGNAL SIGINT
CMD ./mooney serve
