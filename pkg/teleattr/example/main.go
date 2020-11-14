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
