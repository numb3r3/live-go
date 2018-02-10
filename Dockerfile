FROM golang:latest
MAINTAINER Numb3r3 "wangfelix98@gmail.com"

# Copy the directory into the container.
RUN mkdir -p /go/src/github.com/numb3r3/h5-rtms-server/
WORKDIR /go/src/github.com/numb3r3/h5-rtms-server/
ADD . /go/src/github.com/numb3r3/h5-rtms-server/

# Download and install any required third party dependencies into the container.
RUN apk add --no-cache g++ \
&& go-wrapper install \
&& apk del g++

# Expose emitter ports
EXPOSE 4000
EXPOSE 8080
EXPOSE 8443

# Start the broker
CMD ["go-wrapper", "run"]