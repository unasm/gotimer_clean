package test

import (
	"encoding/json"
	"io"
	"net/http"
)

type ret struct {
	ReturnCode        int    `json:"returnCode"`
	ReturnMessage     string `json:"returnMessage"`
	ReturnUserMessage string `json:"returnUserMessage"`
}

type OutData struct {
	Data  string `json:"data"`
	Count int    `json:"count"`
}
type res struct {
	Error ret `json:"error"`
	//Error ret    `json:"error"`
	//Data string `json:"data"`
	Data OutData `json:"data"`
}

func nake(w http.ResponseWriter, r *http.Request) {
	mystruct := res{
		Error: ret{
			ReturnCode:        200,
			ReturnMessage:     "ok",
			ReturnUserMessage: "212",
		},
		//Data: "hello",
		Data: OutData{
			Data:  "hello",
			Count: 2,
		},
	}
	json, _ := json.Marshal(mystruct)
	io.WriteString(w, string(json))
}
