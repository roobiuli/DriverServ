FROM golang:alpine


# Service Label for Datadog APM
LABEL com.datadoghq.tags.env="Testing"
LABEL com.datadoghq.tags.service="Drivers"
LABEL com.datadoghq.tags.version="1"

## Install GO
#RUN apk add --no-cache git make musl-dev go


# Configure Go
#ENV GOROOT /usr/lib/go
#ENV GOPATH /go
#ENV PATH /go/bin:$PATH

# Prepare SourceCode
RUN mkdir /Proj /var/log/applogs/
#COPY ./ /Proj
#RUN cd /Proj ; go mod vendor

WORKDIR /Proj/DriverService/

CMD ["go", "run", "./main.go"]