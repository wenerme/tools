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
