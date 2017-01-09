package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"security/lib/check"
)

var isMock = false

func SetMock() {
	isMock = true
}

/*
	向其他系统发http请求
*/
/*
func HttpPost(url string, params map[string]string) string {
	req := httplib.Post(url)
	//req.Header("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range params {
		req.Param(k, v)
	}
	//req.Debug(true)
	//fmt.Println(req)
	res, err := req.String()
	fmt.Println("updateing")
	fmt.Println(res)
	fmt.Println()
	check.Err(err)
	return res
}
*/

/*
	向其他系统发http请求
*/
func HttpPost(urlStr string, params map[string]string) string {
	postValue := url.Values{}
	for k, v := range params {
		postValue.Add(k, v)
	}
	resp, err := http.PostForm(urlStr, postValue)
	if isMock {
		return `{
			"error": {
			    "returnCode": 0,
				"returnMessage": "success",
				"returnUserMessage": "success"
			},
			"data": 
				[
				"127.0.0.1"
				]
		}`
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if resp == nil || resp.StatusCode != 200 {
		check.Err(fmt.Errorf("调用waf系统异常"))
	}
	body, err := ioutil.ReadAll(resp.Body)
	check.Err(err)
	return string(body)
}

/*
	向其他系统发http请求
*/
/*
func HttpPost(urlStr string, params map[string]string) string {

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader("name=cjb"))
	check.Err(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println("bodying")
	fmt.Println(resp.StatusCode)
	fmt.Println(params)
	fmt.Println(urlStr)
	fmt.Println(string(body))
	fmt.Println()
	fmt.Println()
	return string(body)
}
*/
