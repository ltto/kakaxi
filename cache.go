package kakaxi

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func SaveCache(savePath, host, hostPath string, header http.Header, body []byte) {
	var html = false
	if hs, ok := header["Content-Type"]; ok {
		for _, h := range hs {
			if html = strings.Contains(h, "html"); html {
				break
			}
		}
	}
	//host := URL.Host + URL.Path
	if html {
		body, _ = ReplaceHTML(body, hostPath, host)
	}
	_ = ioutil.WriteFile(savePath, body, 0777)
}
