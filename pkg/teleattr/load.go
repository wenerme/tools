package teleattr

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"
)

func LoadFile(f string) (*PhoneData, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return LoadBytes(b)
}

func LoadBytes(b []byte) (*PhoneData, error) {
	/*
	   | 4 bytes |                     <- phone.dat 版本号（如：1701即17年1月份）
	   ------------
	   | 4 bytes |                     <-  第一个索引的偏移
	   -----------------------
	   |  offset - 8            |      <-  记录区 - <省份>|<城市>|<邮编>|<长途区号>\0 - 山东|济南|250000|0531
	   -----------------------
	   |  index                 |      <-  索引区 - <手机号前七位><记录区的偏移><卡类型>
	   -----------------------
	*/
	l := len(b)
	i := 4
	ind := int(binary.LittleEndian.Uint32(b[i:]))
	i += 4
	s := ""
	indexedRecords := make(map[int]*Record)
	var records []Record
	for i < ind-8 {
		id := i
		s, i = str(b, i)
		rec := parseRec(s)
		rec.Offset = id
		records = append(records, rec)
		indexedRecords[id] = &records[len(records)-1]
	}

	idxs := make(map[int]Index)
	var idxes []Index
	for i = ind; i < l-8; i += 9 {
		idx := parseIdx(b[i:])
		idx.Record = indexedRecords[idx.RecordIndex]
		idxs[i] = idx
		idxes = append(idxes, idx)
	}

	data := &PhoneData{}
	data.Version = string(b[:4])
	data.Index = idxes
	data.Records = records
	return data, nil
}

func parseIdx(b []byte) Index {
	return Index{
		Prefix:      int(binary.LittleEndian.Uint32(b)),
		RecordIndex: int(binary.LittleEndian.Uint32(b[4:])),
		Vendor:      VendorType(b[8]),
	}
}
func str(b []byte, i int) (s string, n int) {
	n = i
	for ; b[n] != 0; n++ {

	}

	s = string(b[i:n])
	n++
	return
}
func parseRec(s string) Record {
	// 山东|济南|250000|0531
	split := strings.Split(s, "|")
	return Record{
		Prov: split[0],
		City: split[1],
		Zip:  split[2],
		Zone: split[3],
	}
}

type Index struct {
	Prefix      int
	RecordIndex int
	Record      *Record
	Vendor      VendorType
}

type VendorType byte

func (idx Index) String() string {
	return fmt.Sprintf("<%v><%v><%v>", idx.Prefix, idx.RecordIndex, idx.Vendor)
}

type Record struct {
	Zip    string
	Zone   string
	Prov   string
	City   string
	Offset int
}

func (r Record) String() string {
	return fmt.Sprintf("%v|%v|%v|%v", r.Prov, r.City, r.Zip, r.Zone)
}
func (b VendorType) String() string {
	switch b {
	case CMCC:
		return "中国移动"
	case CMCCV:
		return "移动虚拟运营商"
	case CUCC:
		return "中国联通"
	case CUCCV:
		return "联通虚拟运营商"
	case CTCC:
		return "中国电信"
	case CTCCV:
		return "电信虚拟运营商"
	default:
		return fmt.Sprintf("未知运营商(%v)", byte(b))
	}
}

const (
	CMCC  VendorType = iota + 0x01 // 中国移动
	CUCC                           // 中国联通
	CTCC                           // 中国电信
	CTCCV                          // 电信虚拟运营商
	CUCCV                          // 联通虚拟运营商
	CMCCV                          // 移动虚拟运营商
)
