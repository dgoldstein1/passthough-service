FROM golang:latest

ENV PRJ_DIR $GOPATH/src/github.com/dgoldstein1/passthrough-service
# create project dir
RUN mkdir -p $PRJ_DIR
# add src, service communication ,and docs
COPY . $PRJ_DIR
RUN mkdir -p mkdir -p /opt/services/passthrough-service
COPY ./Gopkg.toml $PRJ_DIR
COPY ./Gopkg.lock $PRJ_DIR

# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH

# install utils
RUN go get github.com/golang/dep/cmd/dep

# copy over source code
WORKDIR $PRJ_DIR

# install dependencies
RUN dep ensure -v

# configure entrypoint
RUN go build
ENV PORT 3000
ENV MESH_ID "NOT_SET"
ENV SERVER_CERT ""
ENV SERVER_KET ""
ENV SERVER_CA ""
ENV USE_TLS "false"

ENTRYPOINT ["./passthrough-service"]

# expose service ports
EXPOSE 10000
EXPOSE 10001
EXPOSE 3000
