# tools
⚙🔩🔧


Build  | Coverage
-------|----
![test and build](https://github.com/wenerme/tools/workflows/test%20and%20build/badge.svg) | [![Coverage Status](https://coveralls.io/repos/github/wenerme/tools/badge.svg?branch=master)](https://coveralls.io/github/wenerme/tools?branch=master)

> 💡
>
> Pre-build binaries can download from Actions artifacts

## apki
* AlpineLinux package indexer
* `pkg/apki`
* [alpine.ink](https://alpine.ink)

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

__Library Usage__

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

## teleattr
* Telephone Number Attribution
* 电话号码归属地查询
* Data source [xluohome/phonedata](https://github.com/xluohome/phonedata)
  * Source code under GPL-3.0

__Library Usage__


```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	teleattrdata "github.com/wenerme/tools/pkg/teleattr/data"
)

func main() {
	data, err := teleattrdata.PhoneData()
	if err != nil {
		panic(err)
	}
	idx, err := data.Search("13565961")
	if err != nil {
		panic(err)
	}
	spew.Dump(idx.Vendor, idx.Record)
	// (teleattr.VendorType) 中国移动
	// (*teleattr.Record)(0xc00013d9b8)(新疆|乌鲁木齐|830000|0991)
}
```

## scel
* 搜狗 scel 词库

```go
package main

import (
	"fmt"
	"github.com/wenerme/tools/pkg/scel"
)

func main() {
	s, err := scel.LoadFile("pkg/scel/testdata/全国省市区县地名.scel", &scel.LoadOptions{
		ResolveWordPinyin: true,
		SkipExt:           true,
	})
	if err != nil {
		panic(err)
	}
	w, err := s.Search("shang", "hai")
	if err != nil {
		panic(err)
	}
	fmt.Println(w.Words[0])
	// 上海
}
```
