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

func SaveCache(URL url.URL, header http.Header, body []byte) {
	var html = false
	if hs, ok := header["Content-Type"]; ok {
		for _, h := range hs {
			if html = strings.Contains(h, "html"); html {
				break
			}
		}
	}
	host := URL.Host + URL.Path

	if html {
		body, _ = ReplaceHTML(body, "http://"+host, "http://"+host)
	}
	if path.Ext(host) == "" {
		host = host + "/index.html"
	}
	host = path.Join("dao/target", host)
	_ = os.MkdirAll(path.Dir(host), 0777)
	_ = ioutil.WriteFile(host, body, 0777)
	marshal, _ := json.Marshal(header)
	_ = ioutil.WriteFile(host+".meta.json", marshal, 0777)
	fmt.Println("数据大小:", len(body), "抓取连接:", host)
}
