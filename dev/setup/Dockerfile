# This Dockerfile is not for use so much as to verify a working dev environment
# from scratch. Building this with `make test-env` will show missing dependencies
# and help clarify what software is necessary for a new development setup.

FROM alpine:3.4
RUN apk --update add go curl make bash git ca-certificates
ENV GOPATH /usr/local
RUN curl https://glide.sh/get | sh
COPY . /usr/local/src/github.com/gliderlabs/gliderlabs.io
WORKDIR /usr/local/src/github.com/gliderlabs/gliderlabs.io
RUN TIMEOUT=true make clobber setup linux dev
