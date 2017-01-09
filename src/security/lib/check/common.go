package check

import (
	"fmt"
	"regexp"
)

func IsIp(ip string) bool {
	pattern := `^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`
	//pattern := `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
	reg, err := regexp.MatchString(pattern, ip)
	Err(err)
	return reg
}

// 判断是数字
func IsNum(num string) bool {
	pattern := `^\d+$`
	//pattern := `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
	reg, err := regexp.MatchString(pattern, num)
	Err(err)
	return reg
}

//判断是否是uri
func IsUri(uri string) bool {
	pattern := `^/([\da-zA-Z\/\-\_]+)$`
	//[^/]+/)
	reg, err := regexp.MatchString(pattern, uri)
	fmt.Println(reg)
	fmt.Println(uri)
	fmt.Println(pattern)
	Err(err)
	return reg
}
