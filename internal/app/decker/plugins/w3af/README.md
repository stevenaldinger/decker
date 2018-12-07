# w3af decker plugin

[w3af: web application attack and audit framework](https://github.com/andresriancho/w3af/)

https://github.com/andresriancho/w3af/blob/master/scripts/all.w3af

http://docs.w3af.org/en/latest/scripts.html

```
# all usage demo!

plugins
output console,text_file
output
output config text_file
set output_file output-w3af.txt
set verbose True
back
output config console
set verbose False
back

crawl all, !bing_spider, !google_spider, !spider_man
crawl

grep all
grep

audit all
audit

bruteforce all
bruteforce

back

target
set target http://moth/w3af/
back

start
```
