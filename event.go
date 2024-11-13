package kakaxi

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/andybalholm/brotli"
)

var OnRequest = func(r Request) {}
var OnResponse = func(r Request, resp Response, doCache bool, savePath, host, hostPath string) {
	var body = resp.Body
	header := resp.Header
	if encode, ok := header["Content-Encoding"]; ok && len(encode) > 0 {
		//delete(header, "Content-Encoding")
		buffer := bytes.NewBuffer(body)
		var bodyR io.Reader
		var err error
		switch strings.ToLower(encode[0]) {
		case "zlib", "deflate":
			if bodyR, err = zlib.NewReader(buffer); err != nil {
				log.Printf("zlib.NewReader err:%v\n", err)
				return
			}
		case "gzip":
			if bodyR, err = gzip.NewReader(buffer); err != nil {
				log.Printf("gzip.NewReader err:%v\n", err)
				return
			}
		case "br":
			// 解码 Brotli 压缩的内容
			bodyR = brotli.NewReader(buffer)
		case "identity":

		default:
			log.Println("unknown Content-Encoding: " + encode[0])
		}
		body, err = ioutil.ReadAll(bodyR)
		if err != nil {
			log.Printf("ioutil.ReadAll err:%v\n", err)
			return
		}
		if len(body) >= 3 {
			if body[0] == 239 && body[1] == 187 && body[2] == 191 {
				body = body[3:]
			}
		}
	}
	if doCache {
		if resp.StatusCode < 400 {
			SaveCache(savePath, host, hostPath, header, body)
			fmt.Println("数据大小:", FormatBytes(float64(len(body))), "抓取连接:", r.URL.Scheme+"://"+r.URL.Host+r.URL.Path)
		}
	}
}
