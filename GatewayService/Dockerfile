FROM golang:alpine

# Service Label for Datadog APM
LABEL com.datadoghq.tags.service="Gateway"
LABEL com.datadoghq.tags.version="1"

ENV DEBIAN_FRONTEND=noninteractive


# Install GO
RUN apk add --no-cache git make musl-dev
#RUN apt-get update && apt-get install -y git musl-dev vim
# Configure Go
#ENV GOROOT /usr/lib/go
#ENV GOPATH /go
#ENV PATH /go/bin:$PATH
# Prepare Project
#RUN export PATH=/usr/lib/go-1.10/bin:$PATH
RUN mkdir /Proj /var/log/applogs/
#COPY ./ /Proj
#RUN cd /Proj ; go mod vendor

WORKDIR /Proj/GatewayService/

CMD ["go", "run", "./main.go"]



