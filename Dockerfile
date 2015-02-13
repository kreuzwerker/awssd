FROM odise/busybox-curl

MAINTAINER joern.barthel@kreuzwerker.de

ENV SERIAL 20150213

RUN curl -sLo /awssd https://github.com/kreuzwerker/awssd/releases/download/v0.0.1/awssd-linux && chmod +x /awssd

ENTRYPOINT [ "/awssd" ]
