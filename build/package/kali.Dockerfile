FROM stevenaldinger/decker:latest as decker

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
 # && gem install bundler \
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

# add decker to the path
ENV PATH="$PATH:/go/bin"
# decker expects this to exist for the reports it generates
RUN mkdir -p /tmp/reports

COPY --from=decker /go/bin/decker /go/bin/decker
COPY --from=decker /go/bin/internal/app/decker/plugins /go/bin/internal/app/decker/plugins
COPY --from=decker /go/bin/examples /go/bin/examples

# ENTRYPOINT ["bash"]
CMD ["bash"]
