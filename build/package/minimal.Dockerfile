FROM golang:1.12.1 as builder

# install dep
RUN apt-get update \
 && apt-get install -y \
      ca-certificates \
 && curl https://raw.githubusercontent.com/golang/dep/master/install.sh \
      --output /tmp/install-dep.sh \
      --silent \
 && chmod a+x /tmp/install-dep.sh \
 && /tmp/install-dep.sh \
 && rm /tmp/install-dep.sh \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ARG APP_NAME=decker

WORKDIR /go/src/github.com/stevenaldinger/$APP_NAME

COPY . .

RUN dep ensure -v \
 && make build_all \
 && chmod a+x ./$APP_NAME

FROM scratch as decker

# add decker to the path
ENV PATH="$PATH:/go/bin"

COPY --from=builder /go/src/github.com/stevenaldinger/decker/decker /go/bin/decker
COPY --from=builder /go/src/github.com/stevenaldinger/decker/internal/app/decker/plugins /go/bin/internal/app/decker/plugins
COPY --from=builder /go/src/github.com/stevenaldinger/decker/examples /go/bin/examples

# decker expects this to exist for the reports it generates
# RUN mkdir -p /tmp/reports

CMD ["decker"]
