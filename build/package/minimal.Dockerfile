FROM golang:1.12.1 as builder

ARG UNAME_S="Linux"
ARG GOARCH="amd64"
ARG TEMP_DL_DIR="/tmp/downloads"
ARG GOLANG_DL_BASE_URL="https://dl.google.com/go/go"
ARG GO_DEP_INSTALL_SCRIPT="https://raw.githubusercontent.com/golang/dep/master/install.sh"
ARG GO_DEP_RELEASE_TAG="v0.5.1"
ARG YOUR_GITHUB_HANDLE="stevenaldinger"
ARG APP_NAME="decker"

ENV \
 DEP_RELEASE_TAG="${GO_DEP_RELEASE_TAG}" \
 GOPATH="/go" \
 GOBIN="/go/bin" \
 GOARCH="amd64" \
 PATH="$PATH:/go/bin"

# Install golang and dep
# Find versions and DL links here: https://golang.org/dl/
RUN apt-get update \
 && apt-get -y install \
    build-essential \
    ca-certificates \
    curl \
    git \
 && mkdir -p "${GOBIN}" "${TEMP_DL_DIR}" \
 && go version \
 && curl "${GO_DEP_INSTALL_SCRIPT}" \
     --output "${TEMP_DL_DIR}/install-dep.sh" \
     --silent \
 && chmod a+x "${TEMP_DL_DIR}/install-dep.sh" \
 && cat "${TEMP_DL_DIR}/install-dep.sh" \
 && "${TEMP_DL_DIR}/install-dep.sh" \
 && rm "${TEMP_DL_DIR}/install-dep.sh" \
 && go get -u golang.org/x/lint/golint \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* "${TEMP_DL_DIR}"

# # install dep
# RUN apt-get update \
#  && apt-get install -y \
#       ca-certificates \
#  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh \
#       --output /tmp/install-dep.sh \
#       --silent \
#  && chmod a+x /tmp/install-dep.sh \
#  && /tmp/install-dep.sh \
#  && rm /tmp/install-dep.sh \
#  && apt-get clean \
#  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ARG APP_NAME=decker

WORKDIR /go/src/github.com/stevenaldinger/$APP_NAME

COPY . .

RUN dep ensure -v \
 && make build_all \
 && chmod a+x ./$APP_NAME

FROM scratch as decker

COPY --from=builder /go/src/github.com/stevenaldinger/decker/decker /go/bin/decker
COPY --from=builder /go/src/github.com/stevenaldinger/decker/internal/app/decker/plugins /go/bin/internal/app/decker/plugins
COPY --from=builder /go/src/github.com/stevenaldinger/decker/examples /go/bin/examples

# decker expects this to exist for the reports it generates
# RUN mkdir -p /tmp/reports

CMD ["decker"]
