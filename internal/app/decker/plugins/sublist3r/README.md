# Sublist3r decker plugin

[Fast subdomains enumeration tool for penetration testers](https://github.com/aboul3la/Sublist3r)

## Installation

```sh
# Sublist3r
RUN apt-get update \
 && apt-get install -y \
      python-3 \
      python3-pip \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
 && git clone https://github.com/aboul3la/Sublist3r.git /usr/bin/Sublist3r \
 && pip3 install -r /usr/bin/Sublist3r/requirements.txt
```
