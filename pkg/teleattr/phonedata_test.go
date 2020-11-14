package teleattr_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/wenerme/tools/pkg/teleattr"
	teleattrdata "github.com/wenerme/tools/pkg/teleattr/data"
)

func TestLoad(t *testing.T) {
	d, err := teleattr.LoadFile("./testdata/phone.dat")
	if !assert.NoError(t, err) {
		return
	}
	check(t, d)
}

func check(t *testing.T, d *teleattr.PhoneData) {
	{
		idx, err := d.Search("185215912")
		if assert.NoError(t, err) {
			spew.Dump(idx, idx.Record)
			assert.Equal(t, "上海", idx.Record.City)
			assert.Equal(t, "中国联通", idx.Vendor.String())
		}
	}
	{
		_, err := d.Search("1852")
		assert.Equal(t, err, teleattr.ErrInvalidNumber)
	}
	{
		_, err := d.Search("9999999999")
		assert.Equal(t, err, teleattr.ErrNotFound)
	}
}

func TestStaticData(t *testing.T) {
	d, err := teleattrdata.PhoneData()
	if !assert.NoError(t, err) {
		return
	}
	check(t, d)
}
