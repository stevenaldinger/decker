FROM kalilinux/kali-linux-docker:latest as development

ARG GO_INSTALL_DIR_PREFIX="/usr/local"
ARG GOPATH="/go"
ARG GOBIN="${GOPATH}/bin"
ARG UNAME_S="Linux"
ARG GOARCH="amd64"
ARG TEMP_DL_DIR="/tmp/downloads"
ARG GOLANG_DL_BASE_URL="https://dl.google.com/go/go"
ARG GO_DEP_INSTALL_SCRIPT="https://raw.githubusercontent.com/golang/dep/master/install.sh"
ARG GO_DEP_RELEASE_TAG="v0.5.1"

ARG GOLANG_VERSION="1.12.1"
ARG GO_CHECKSUM="2a3fdabf665496a0db5f41ec6af7a9b15a49fbe71a85a50ca38b1f13a103aeec"
ARG YOUR_GITHUB_HANDLE="stevenaldinger"
ARG APP_NAME="decker"

ENV \
 GOLANG_VERSION="${GOLANG_VERSION}" \
 DEP_RELEASE_TAG="${GO_DEP_RELEASE_TAG}" \
 GOPATH="/go" \
 GOBIN="/go/bin" \
 GOARCH="amd64" \
 GOROOT="${GO_INSTALL_DIR_PREFIX}/go" \
 GO_BIN_PATH_HOST="${TEMP_DL_DIR}/go${GOLANG_VERSION}.${UNAME_S}-${GOARCH}.tar.gz" \
 PATH="$PATH:/go/bin:${GO_INSTALL_DIR_PREFIX}/go/bin"

ENV GO_BIN_URL_REMOTE="${GOLANG_DL_BASE_URL}${GOLANG_VERSION}.${UNAME_S}-${GOARCH}.tar.gz"

# Install golang and dep
# Find versions and DL links here: https://golang.org/dl/
RUN apt-get update \
 && apt-get -y install \
     build-essential \
     ca-certificates \
     curl \
     git \
 && mkdir -p "${GOBIN}" "${TEMP_DL_DIR}" "${GO_INSTALL_DIR_PREFIX}/go" \
 && curl -L "${GO_BIN_URL_REMOTE}" \
      --output "${GO_BIN_PATH_HOST}" \
      --silent \
 && tar -C "${GO_INSTALL_DIR_PREFIX}" -zxf "${GO_BIN_PATH_HOST}" \
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

# install ruby, bundler, nokogiri(necessary for bundler install)
# zlib (necessary for nokogiri install)
# install wpscan and nmap as decker dependencies at runtime
# dnsutils - nslookup, dig, host
RUN apt-get update \
 && apt-get install -y \
      ca-certificates \
      dnsutils \
      git \
      nmap \
      python3 \
      python3-pip \
      python-pip \
      ruby-full \
      sslscan \
      whois \
      zlib1g-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
 # && gem install nokogiri -v '1.8.5' --source 'https://rubygems.org/' \
 # && gem install bundler -v '1.17' \
 # && git clone https://github.com/wpscanteam/wpscan /usr/bin/wpscan \
 # && cd /usr/bin/wpscan/ \
 # && bundle install \
 # && rake install \
 && pip install --upgrade setuptools \
 && pip install --upgrade sslyze \
 && pip3 install --upgrade setuptools \
 && git clone https://github.com/laramies/theHarvester.git /usr/bin/theHarvester \
 && pip3 install -r /usr/bin/theHarvester/requirements.txt \
 && git clone https://github.com/darkoperator/dnsrecon.git /usr/bin/dnsrecon \
 && pip install -r /usr/bin/dnsrecon/requirements.txt \
 && git clone https://github.com/aboul3la/Sublist3r.git /usr/bin/Sublist3r \
 && pip3 install -r /usr/bin/Sublist3r/requirements.txt \
 && git clone https://github.com/EnableSecurity/wafw00f.git /usr/bin/wafw00f \
 && cd /usr/bin/wafw00f \
 && python setup.py install \
 && git clone https://github.com/khalilbijjou/WAFNinja /usr/bin/WAFNinja \
 && pip install -r /usr/bin/WAFNinja/requirements.txt \
 && git clone https://github.com/faizann24/XssPy.git /usr/bin/XssPy \
 && pip install mechanize \
 && git clone https://github.com/hahwul/a2sv.git /usr/bin/a2sv \
 && pip install -r /usr/bin/a2sv/requirements.txt

RUN apt-get update \
 && apt-get install -y kali-linux-full

WORKDIR /go/src/github.com/${YOUR_GITHUB_HANDLE}/${APP_NAME}

COPY . .

RUN dep ensure -v \
 && make build_all \
 && chmod a+x ./$APP_NAME

CMD ["bash"]
