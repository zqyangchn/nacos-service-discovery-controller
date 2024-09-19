FROM ubuntu:22.04

LABEL maintainer="zqyangchn@gmail.com" description="nacos-service-discovery-controller"

RUN apt-get -qq update && \
        export DEBIAN_FRONTEND=noninteractive; \
        export DEBCONF_NONINTERACTIVE_SEEN=true; \
        echo 'tzdata tzdata/Areas select Asia' | debconf-set-selections; \
        echo 'tzdata tzdata/Zones/Asia select Shanghai' | debconf-set-selections; \
        apt update -qqy && apt -qqy upgrade && \
        apt install -qqy --no-install-recommends tzdata && \
        apt install -qqy procps iproute2 net-tools iputils-ping telnet \
            sudo libaio1 libfontconfig1 libxrender1 libxext6 fontconfig \
            htop curl tcpdump lsof vim sysstat zip openssh-client && \
        apt autoclean && rm -rf /var/lib/apt/lists/* && fc-cache -fv

ADD nacos-service-discovery-controller /usr/bin/nacos-service-discovery-controller

EXPOSE 8428
ENTRYPOINT ["/usr/bin/nacos-service-discovery-controller"]
