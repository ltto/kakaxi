package kakaxi

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"strings"
)

var OnRequest = func(r Request) {}
var OnResponse = func(r Request, resp Response) {
	var body = resp.Body
	header := resp.Header
	if encode, ok := header["Content-Encoding"]; ok && len(encode) > 0 {
		//delete(header, "Content-Encoding")
		buffer := bytes.NewBuffer(body)
		var bodyR io.ReadCloser
		var err error
		switch strings.ToLower(encode[0]) {
		case "zlib", "deflate":
			if bodyR, err = zlib.NewReader(buffer); err != nil {
				return
			}
		case "gzip":
			if bodyR, err = gzip.NewReader(buffer); err != nil {
				return
			}
		case "identity":
		}
		body, err = ioutil.ReadAll(bodyR)
		if err != nil {
			panic(err)
		}
	}
	SaveCache(*r.URL, header, body)
}
