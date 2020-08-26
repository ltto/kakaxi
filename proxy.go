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

type Buffer struct {
	*bytes.Buffer
}

func (b Buffer) Close() error {
	return nil
}

func ProxyHTTP(r http.Request) (headerD http.Header, resp *http.Response, bodyB []byte, err error) {
	r.RequestURI = ""
	if r.URL.Host == "" {
		r.URL.Host = r.Host
		r.URL.Scheme = "https"
	}
	fmt.Println("host-----:", r.Host)
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}
	if OnRequest != nil {
		var on = r
		on.Body = Buffer{bytes.NewBuffer(body)}
		OnRequest(on)
	}

	resp, err = http.DefaultClient.Do(&r)
	if err != nil {
		return
	}
	var bodyR = resp.Body
	headerD = resp.Header

	if encode, ok := resp.Header["Content-Encoding"]; ok && len(encode) > 0 {
		delete(headerD, "Content-Encoding")
		switch strings.ToLower(encode[0]) {
		case "zlib", "deflate":
			if bodyR, err = zlib.NewReader(resp.Body); err != nil {
				return
			}
		case "gzip":
			if bodyR, err = gzip.NewReader(resp.Body); err != nil {
				return
			}
		case "identity":
		}
	}
	bodyB, err = ioutil.ReadAll(bodyR)
	if OnResponse != nil {
		var on = r
		on.Body = Buffer{bytes.NewBuffer(body)}
		OnRequest(on)
		OnResponse(on, headerD, bodyB)
	}
	return
}

func Writer(w io.Writer, header http.Header, resp *http.Response, body []byte) {
	_, _ = w.Write([]byte(fmt.Sprintf("%s %s\n", resp.Proto, resp.Status)))
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
