FROM golang:1.12.1

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

# install ruby, bundler, nokogiri(necessary for bundler install)
# zlib (necessary for nokogiri install)
# install wpscan and nmap as decker dependencies at runtime
# dnsutils - nslookup, dig, host
RUN apt-get update \
 && apt-get install -y \
      dnsutils \
      nmap \
      ruby-full \
      zlib1g-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
 && gem install nokogiri -v '1.8.5' --source 'https://rubygems.org/' \
 && gem install bundler -v '1.17' \
 && git clone https://github.com/wpscanteam/wpscan /usr/bin/wpscan \
 && cd /usr/bin/wpscan/ \
 && bundle install \
 && rake install

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y \
      net-tools \
      python-pip \
      software-properties-common \
      whois \
      python3

# These errors happen with other python-3.6 repos but things seem to work alright
# so just ignore on repo add/remove with ";" instead of using "&&"
# E: Failed to fetch http://ppa.launchpad.net/jonathonf/python-3.6/ubuntu/dists/disco/main/binary-amd64/Packages  404  Not Found
# E: Some index files failed to download. They have been ignored, or old ones used instead.
RUN apt-get update \
 # these make sure pip3 gets built right
 && apt-get install -y \
      libreadline-gplv2-dev \
      libncursesw5-dev \
      libssl-dev \
      libsqlite3-dev \
      tk-dev \
      libgdbm-dev \
      libc6-dev \
      libbz2-dev \
 && cd /opt \
 && wget https://www.python.org/ftp/python/3.6.7/Python-3.6.7.tgz \
 && tar -xvf Python-3.6.7.tgz \
 && cd Python-3.6.7 \
 && ./configure \
 && make \
 && make install \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
 && pip install --upgrade setuptools \
 && pip install --upgrade sslyze \
 && pip3 install --upgrade setuptools

# RUN apt-get update \
#  && apt-get install -y \
#       python3-pip \
#  && apt-get clean \
#  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
#  && pip install --upgrade setuptools \
#  && pip install --upgrade sslyze \
#  && pip3 install --upgrade setuptools

# whois - whois plugin
# python-pip - sslyze plugin
# RUN apt-get update \
#  && apt-get install -y \
#       python3 \
#       python3-pip \
#       python-pip \
#       whois \
#  && apt-get clean \
#  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
#  && pip install --upgrade setuptools \
#  && pip install --upgrade sslyze \
#  && pip3 install --upgrade setuptools

# The harvester download and dependencies
RUN git clone https://github.com/laramies/theHarvester.git /usr/bin/theHarvester \
 && pip3 install -r /usr/bin/theHarvester/requirements.txt

# DNSRecon download and Dependencies
RUN git clone https://github.com/darkoperator/dnsrecon.git /usr/bin/dnsrecon \
 && pip install -r /usr/bin/dnsrecon/requirements.txt

# Sublist3r
RUN git clone https://github.com/aboul3la/Sublist3r.git /usr/bin/Sublist3r \
 && pip3 install -r /usr/bin/Sublist3r/requirements.txt

# wafw00f
RUN git clone https://github.com/EnableSecurity/wafw00f.git /usr/bin/wafw00f \
 && cd /usr/bin/wafw00f \
 && python setup.py install

# WAFNinja
RUN git clone https://github.com/khalilbijjou/WAFNinja /usr/bin/WAFNinja \
 && pip install -r /usr/bin/WAFNinja/requirements.txt

# XssPy (discovered issue with the Pentagon)
RUN git clone https://github.com/faizann24/XssPy.git /usr/bin/XssPy \
 && pip install mechanize

# a2sv - auto scanning to SSL vulnerability
RUN git clone https://github.com/hahwul/a2sv.git /usr/bin/a2sv \
 && pip install -r /usr/bin/a2sv/requirements.txt

# sslscan
RUN apt-get update \
 && apt-get install -y \
      sslscan \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ARG APP_NAME=decker

WORKDIR /go/src/github.com/stevenaldinger/$APP_NAME

RUN go get -u golang.org/x/lint/golint

RUN echo "Backing up /etc/apt/sources.list..." \
 && cp -a /etc/apt/sources.list /etc/apt/sources.list.bak \
 && echo "Adding Kali sources to /etc/apt/sources.list..." \
 && echo "deb http://http.kali.org/kali kali-rolling main contrib non-free" >> /etc/apt/sources.list \
 && echo "deb-src http://http.kali.org/kali kali-rolling main contrib non-free" >> /etc/apt/sources.list \
 && apt-get update \
 && apt-get install -y metasploit-framework --allow-unauthenticated \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
 && echo "Restoring /etc/apt/sources.list..." \
 && rm /etc/apt/sources.list \
 && mv /etc/apt/sources.list.bak /etc/apt/sources.list

RUN git clone --single-branch -b feature/support-noninteractive https://www.github.com/stevenaldinger/routersploit /usr/bin/routersploit \
 && cd /usr/bin/routersploit \
 && python3 -m pip install -r requirements.txt

# doesn't work without better wireless adapter
# RUN apt-get update \
#  && apt-get install -y dnsmasq \
#  && git clone https://github.com/wifiphisher/wifiphisher.git /usr/bin/wifiphisher \
#  && cd /usr/bin/wifiphisher \
#  && python setup.py install

RUN apt-get update \
 && apt-get install -y \
      pciutils \
      aircrack-ng \
      kmod \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN echo "Backing up /etc/apt/sources.list..." \
 && cp -a /etc/apt/sources.list /etc/apt/sources.list.bak \
 && echo "Adding Kali sources to /etc/apt/sources.list..." \
 && echo "deb http://http.kali.org/kali kali-rolling main contrib non-free" >> /etc/apt/sources.list \
 && echo "deb-src http://http.kali.org/kali kali-rolling main contrib non-free" >> /etc/apt/sources.list \
 && apt-get update \
 && apt-get install -y bluelog --allow-unauthenticated \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
 && echo "Restoring /etc/apt/sources.list..." \
 && rm /etc/apt/sources.list \
 && mv /etc/apt/sources.list.bak /etc/apt/sources.list

# nodejs / npm are requirements for w3af
RUN curl -sL https://deb.nodesource.com/setup_8.x --output /tmp/nodejs-deb \
 && bash /tmp/nodejs-deb \
 && apt-get install -y nodejs

RUN apt-get update \
 && apt-get install -y \
      libxml2-dev \
      libxslt1-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN git clone https://github.com/andresriancho/w3af.git /usr/bin/w3af \
 && cd /usr/bin/w3af \
 && ./w3af_console \
 ; . /tmp/w3af_dependency_install.sh

COPY . .

CMD ["bash"]
