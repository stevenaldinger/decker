FROM stevenaldinger/decker:minimal as decker

# FROM stevenaldinger/w3af:0.0.1 as w3af

FROM kalilinux/kali-linux-docker

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

# RUN curl -Lo /tmp/python-support_1.0.15_all.deb http://ftp.br.debian.org/debian/pool/main/p/python-support/python-support_1.0.15_all.deb \
#  && dpkg -i /tmp/python-support_1.0.15_all.deb \
#  && curl -Lo /tmp/python-webkit_1.1.8-3_amd64.deb http://ftp.br.debian.org/debian/pool/main/p/pywebkitgtk/python-webkit_1.1.8-3_amd64.deb \
#  && dpkg -i /tmp/python-webkit_1.1.8-3_amd64.deb
#
#  && curl -sL https://deb.nodesource.com/setup_8.x --output /tmp/nodejs-deb \
#  && bash /tmp/nodejs-deb \
#  && apt-get install -y nodejs \
#  && npm install -g retire \
#  && git clone https://github.com/andresriancho/w3af.git /usr/bin/w3af \
#  && cd /usr/bin/w3af \
#  && ./w3af_console \
#  ; . /tmp/w3af_dependency_install.sh

# add decker and w3af to the path
ENV PATH="$PATH:/go/bin:/usr/bin/w3af"
# decker expects this to exist for the reports it generates
RUN mkdir -p /tmp/reports

RUN apt-get update \
 && apt-get install -y \
      libxml2-dev \
      libxslt1-dev \
      python-dev \
      zlib1g-dev \
      nodejs \
      npm \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=decker /go/bin/decker /go/bin/decker
COPY --from=decker /go/bin/internal/app/decker/plugins /go/bin/internal/app/decker/plugins
COPY --from=decker /go/bin/examples /go/bin/examples

# COPY --from=w3af /home/w3af/w3af /usr/bin/w3af
#
# ENV LC_ALL=C
#
# # w3af_console fails but craete /tmp/w3af_dependency_install.sh
# # sed changes the install script to add the -y and not require input
# RUN cd /usr/bin/w3af \
#  && pip install --upgrade pip \
#  &&  ./w3af_console \
#  ; sed 's/sudo //g' -i /tmp/w3af_dependency_install.sh \
#  && sed 's/apt-get/apt-get -y/g' -i /tmp/w3af_dependency_install.sh \
#  && sed 's/pip install/pip install --upgrade/g' -i /tmp/w3af_dependency_install.sh \
#  && /tmp/w3af_dependency_install.sh \
#  ; ./w3af_gui \
#  ; sed 's/sudo //g' -i /tmp/w3af_dependency_install.sh \
#  && sed 's/apt-get/apt-get -y/g' -i /tmp/w3af_dependency_install.sh \
#  && sed 's/pip install/pip install --upgrade/g' -i /tmp/w3af_dependency_install.sh \
#  && /tmp/w3af_dependency_install.sh \
#  # Cleanup to make the image smaller
#  ; rm /tmp/w3af_dependency_install.sh \
#  && apt-get clean \
#  && rm -rf /var/lib/apt/lists/* \
#  && rm -rf /tmp/pip-build-root \
#  && sed "s/'accepted-disclaimer': 'false'/'accepted-disclaimer': 'true'/g" -i /usr/bin/w3af/w3af/core/data/db/startup_cfg.py \
#  && sed "s/'skip-dependencies-check': 'false'/'skip-dependencies-check': 'true'/g" -i /usr/bin/w3af/w3af/core/data/db/startup_cfg.py
#
# RUN pip install \
#      acora \
#      bravado_core \
#      diff_match_patch \
#      esmre \
#      pebble \
#      tldextract


# ENTRYPOINT ["bash"]
CMD ["bash"]
