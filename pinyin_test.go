package tests

import (
	"github.com/mozillazg/go-pinyin"
	"testing"
)

func TestPinyin(t *testing.T) {
	args := pinyin.NewArgs()
	args.Style = pinyin.FirstLetter
	args.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	s := "你好.A啊"

	initials := pinyin.Pinyin(s, args)
	t.Log(initials)

	for _, v1 := range initials {
		for _, v2 := range v1 {
			t.Log(v2)
		}
	}
}
