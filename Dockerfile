FROM golang:1.13
MAINTAINER Gursimran Singh <singhgursimran@me.com>

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin
ARG BUILD_ID
ENV BUILD_IMAGE=$BUILD_ID
ENV GO111MODULE=off

# build directories
ADD . /go/src/git.xenonstack.com/xs-onboarding/accounts
WORKDIR /go/src/git.xenonstack.com/xs-onboarding/accounts

# Go dep!
#RUN go get -u github.com/golang/dep/...
#RUN dep ensure -update

RUN go install git.xenonstack.com/xs-onboarding/accounts

EXPOSE 8000

