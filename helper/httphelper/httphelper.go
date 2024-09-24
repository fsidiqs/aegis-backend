package httphelper

import (
	"bytes"
	"net/http"
)

func MakePostJson(url string, body []byte) {
	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"*/*"},
	}
	data := bytes.NewBuffer(body)
	request, _ := http.NewRequest(http.MethodPost, url, data)
	request.Header = headers

	client := &http.Client{}
	client.Do(request)
}
