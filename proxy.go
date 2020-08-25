package kakaxi

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func ProxyHTTP(host, method string, header http.Header, body []byte) (headerD http.Header, bodyB []byte, err error) {
	fmt.Println("host-----:", host)
	if OnRequest != nil {
		OnRequest(header, host, body)
	}

	request, err := http.NewRequest(method, host, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	request.Header = header
	var do *http.Response
	do, err = http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	var bodyR = do.Body
	headerD = do.Header

	if encode, ok := do.Header["Content-Encoding"]; ok && len(encode) > 0 {
		delete(headerD, "Content-Encoding")
		switch strings.ToLower(encode[0]) {
		case "zlib", "deflate":
			if bodyR, err = zlib.NewReader(do.Body); err != nil {
				return
			}
		case "gzip":
			if bodyR, err = gzip.NewReader(do.Body); err != nil {
				return
			}
		case "identity":
		}
	}
	bodyB, err = ioutil.ReadAll(bodyR)
	if OnResponse != nil {
		OnResponse(headerD, host, bodyB)
	}
	return
}

func Writer(w io.Writer, header http.Header, body []byte) {
	_, _ = w.Write([]byte("HTTP/1.1 200 OK\n"))
	header["Content-Length"] = []string{strconv.Itoa(len(body))}
	for k, hs := range header {
		for _, h := range hs {
			if h == "keep-alive" {
				continue
			}
			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s\n", k, h)))
		}
	}
	_, _ = w.Write([]byte("\n"))
	_, _ = w.Write(body)
}
