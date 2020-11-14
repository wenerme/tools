# tools
⚙🔩🔧


Build  | Coverage
-------|----
![test and build](https://github.com/wenerme/tools/workflows/test%20and%20build/badge.svg) | [![Coverage Status](https://coveralls.io/repos/github/wenerme/tools/badge.svg?branch=master)](https://coveralls.io/github/wenerme/tools?branch=master)

## AlpineLinux package indexer
* `pkg/apki`

## libmagic
linmagic Golang wrapper

```go
package main
import (
	"fmt"
	"os"
	"github.com/wenerme/tools/pkg/libmagic"
)

func main() {
    mgc := libmagic.Open(libmagic.MAGIC_NONE)
    defer mgc.Close()
    fmt.Printf("file: %s", mgc.File(os.Args[0]))
    mgc.SetFlags(libmagic.MAGIC_MIME | libmagic.MAGIC_MIME_ENCODING)
    fmt.Printf("file: %s", mgc.File(os.Args[0]))
}
```

## crontimer
Minimal WebCron

> Pre build binaries can download from Actions artifacts

```bash
go get -u github.com/wenerme/tools/cmd/crontimer
# list jobs
crontimer -c cmd/crontimer/crontimer.yaml list
# run cron
crontimer -c cmd/crontimer/crontimer.yaml
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

## stardict
* pkg/stardict
  * reader

```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/wenerme/tools/pkg/stardict"
)

func main() {
	// open dict archive or dir
	dict, err := stardict.Open("./stardict-xhzd-2.4.2.tar.bz2")
	if err != nil {
		panic(err)
	}
	spew.Dump(dict.Info)
	spew.Dump(dict.Search("你好"))
}
```
