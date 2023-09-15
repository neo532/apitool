package http

/*
 * @abstract 传输协议http的一些通用方法
 * @mail neo532@126.com
 * @date 2023-09-12
 */

import (
	"strings"
)

func GetFromCookieString(cookies string, fields ...string) (values []string) {
	lenFields := len(fields)
	ckis := strings.Split(cookies, ";")
	values = make([]string, lenFields)
	var i int
	for _, cookie := range ckis {
		kv := strings.Split(strings.TrimSpace(cookie), "=")
		if i >= lenFields {
			break
		}
		if len(kv) == 2 {
			for offset, f := range fields {
				if f == kv[0] {
					values[offset] = kv[1]
					i++
				}
			}
		}
	}
	return
}

func InCookieString(cookies string, field string) (has bool) {
	return strings.Index(cookies, field+"=") != -1
}
