package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

// GetCityByIp 获取ip所属城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}
	if ip == "::1" || ip == "127.0.0.1" {
		return "内网IP"
	}
	// codecName0 := gjson.Get(s, "streams.0.codec_name").String()

	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	src := string(body)
	srcCharset := "GBK"
	tmp, _ := ToUTF8(srcCharset, src)
	if gjson.Get(tmp, "code").Int() == 0 {
		city := fmt.Sprintf("%s %s", gjson.Get(tmp, "pro").String(), gjson.Get(tmp, "city").String())
		return city
	} else {
		return ""
	}
}

var (
	// Alias for charsets.
	charsetAlias = map[string]string{
		"HZGB2312": "HZ-GB-2312",
		"hzgb2312": "HZ-GB-2312",
		"GB2312":   "HZ-GB-2312",
		"gb2312":   "HZ-GB-2312",
	}
)

// Supported returns whether charset `charset` is supported.
func Supported(charset string) bool {
	return getEncoding(charset) != nil
}

// Convert converts `src` charset encoding from `srcCharset` to `dstCharset`,
// and returns the converted string.
// It returns `src` as `dst` if it fails converting.
func Convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src
	// Converting `src` to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", fmt.Errorf(`convert string "%s" to utf8 failed`, srcCharset)
			}
			src = string(tmp)
		} else {
			return dst, fmt.Errorf(`unsupported srcCharset "%s"`, srcCharset)
		}
	}
	// Do the converting from UTF-8 to `dstCharset`.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", fmt.Errorf(`convert string from utf8 to "%s" failed`, dstCharset)
			}
			dst = string(tmp)
		} else {

			return dst, fmt.Errorf(`unsupported dstCharset "%s"`, dstCharset)
		}
	} else {
		dst = src
	}
	return dst, nil
}

// ToUTF8 converts `src` charset encoding from `srcCharset` to UTF-8 ,
// and returns the converted string.
func ToUTF8(srcCharset string, src string) (dst string, err error) {
	return Convert("UTF-8", srcCharset, src)
}

// UTF8To converts `src` charset encoding from UTF-8 to `dstCharset`,
// and returns the converted string.
func UTF8To(dstCharset string, src string) (dst string, err error) {
	return Convert(dstCharset, "UTF-8", src)
}

// getEncoding returns the encoding.Encoding interface object for `charset`.
// It returns nil if `charset` is not supported.
func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetAlias[charset]; ok {
		charset = c
	}
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		fmt.Errorf(`%+v`, err)
	}
	return enc
}
