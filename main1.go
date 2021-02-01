package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

func UnicodeEmojiCode2(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			continue
		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

func main() {
	str := "6OC0💗💗💗💗"
	fmt.Println(str)

	str1 := UnicodeEmojiDecode(str)
	fmt.Println(str1)

	str2 := UnicodeEmojiCode(str)
	fmt.Println(str2)

	str3 := UnicodeEmojiCode2(str)
	fmt.Println(str3)
}
