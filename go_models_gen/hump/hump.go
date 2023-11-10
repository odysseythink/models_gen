package hump

import (
	"fmt"
	"strings"
)

var LintGonicMapper = []string{
	"API",
	"ASCII",
	"CPU",
	"CSS",
	"DNS",
	"EOF",
	"GUID",
	"HTML",
	"HTTP",
	"HTTPS",
	"ID",
	"IP",
	"JSON",
	"LHS",
	"QPS",
	"RAM",
	"RHS",
	"RPC",
	"SLA",
	"SMTP",
	"SSH",
	"TLS",
	"TTL",
	"UI",
	"UID",
	"UUID",
	"URI",
	"URL",
	"UTF8",
	"VM",
	"XML",
	"XSRF",
	"XSS",
}

// capitalize 字符首字母大写
func capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func BigHumpName(name string) string {
	ret := ""
	tmps := strings.Split(name, "_")
	for _, v := range tmps {
		v1 := strings.ToUpper(v)
		isGonic := false
		for _, v2 := range LintGonicMapper {
			if v1 == v2 {
				isGonic = true
				break
			}
		}
		if isGonic {
			ret += v1
		} else {
			v1 = strings.ToLower(v1)
			v1 = capitalize(v1)
			ret += v1
		}
	}
	return ret
}
