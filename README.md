# tools
⚙🔩🔧

## crontimer
Minimal WebCron

> Pre build binaries can download from Actions artifacts

```bash
go get -u github.com/wenerme/tools/cmd/crontimer
# list jobs
crontimer -c doc/crontimer.yaml list
# run cron
crontimer -c doc/crontimer.yaml
```

```yaml
# cron syntax http://crontab.guru/
jobs:
  - url: http://www.wener.tech
    # interval unit  ms,s,m,h
    interval: 5s
    # or use cron
    # spec: "*/10 * * * *" # 10m
    name: track
  - url: http://www.wener.me
    interval: 5s
    log:
      # log to file
      # set to `.' will use  `job.<index>>.log`
      file: test.log
      # log response
      response: true
log:
  level: debug
  file: log.log
```
