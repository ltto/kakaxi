package kakaxi

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func CopyRequest(request *http.Request) (req Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("CopyRequest:", err)
		return
	}
	temp := *request
	temp.Body = Buffer{bytes.NewBuffer(body)}
	if temp.URL.Host == "" {
		temp.URL.Host = request.Host
		temp.URL.Scheme = "https"
	}
	req.Request = &temp
	req.Body = body
	return
}

type Request struct {
	*http.Request
	Body []byte
}

type Response struct {
	Status     string
	StatusCode int
	Proto      string
	Header     http.Header
	Body       []byte
}

func CopyResponse(response *http.Response) (resp Response) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("CopyResponse:", err)
		return
	}
	resp.StatusCode = response.StatusCode
	resp.Status = response.Status
	resp.Proto = response.Proto
	resp.Body = body
	resp.Header = response.Header
	return
}
