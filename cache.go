package kakaxi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func SaveCache(host string, header http.Header, body []byte) {
	var hostURL, _ = url.Parse(host)
	if hostURL.Path == "" {
		return
	}
	var html = false
	if hs, ok := header["Content-Type"]; ok {
		for _, h := range hs {
			if html = strings.Contains(h, "html"); html {
				break
			}
		}
	}
	if html {
		body, _ = ReplaceHTML(body, "http://demo.qzhai.net/cell/")
	}
	if path.Ext(host) == "" {
		host = host + "/index.html"
	}
	host = path.Join("dao/target", host)
	_ = os.MkdirAll(path.Dir(host), 0777)
	_ = ioutil.WriteFile(host, body, 0777)
	marshal, _ := json.Marshal(header)
	_ = ioutil.WriteFile(host+".meta.json", marshal, 0777)
	fmt.Println(len(body), "抓取连接:", host)
}
