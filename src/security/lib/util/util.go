package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"security/lib/check"
	"strconv"
	"strings"
	"time"
)

//将时间转换为 时间戳
func GetUnix(dateStr string) int64 {
	if len(dateStr) == 0 {
		return 0
	}
	if len(dateStr) <= 10 {
		dateStr = dateStr + " 00:00:00"
	}
	local, err := time.LoadLocation("Asia/Shanghai")
	startTimer, err := time.ParseInLocation("2006-01-02 15:04:05 ", dateStr, local)
	//startTimer, err := time.Parse("2006-01-02 15:04:05 ", dateStr)
	check.Err(err)
	return startTimer.Unix()
}

func GetMin(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func GetMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func IntToStr(data int32) string {
	return ToStr(data)
}

type argInt []int

func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

// ToStr interface to string
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

func GetFile(filePath string) ([]string, error) {
	fp, err := os.Open(filePath)
	check.Err(err)
	buf := bufio.NewReader(fp)
	rsStr := make([]string, 0)
	cnt := 0
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return rsStr, err
		}
		rsStr = append(rsStr, line)
		cnt++
	}
	return rsStr, nil
}

func Implode(strArr *[]string) string {
	buffer := ""
	for _, value := range *strArr {
		buffer = buffer + value + ","
	}
	return strings.Trim(buffer, ",")
}
