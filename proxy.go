package kakaxi

import (
	"fmt"
	"io"
	"net/http"
)

func ProxyHTTP(req Request) (resp Response, err error) {
	if resp, err = doProxyHTTP(req); err != nil {
		return
	}
	return
}
func doProxyHTTP(r Request) (resp Response, err error) {
	r.Request.RequestURI = ""
	var do *http.Response
	do, err = http.DefaultClient.Do(r.Request)
	if err != nil {
		return
	}
	resp = CopyResponse(do)
	return
}

func Writer(w io.Writer, resp Response) {
	_, _ = w.Write([]byte(fmt.Sprintf("%s %s\n", resp.Proto, resp.Status)))
	for k, hs := range resp.Header {
		for _, h := range hs {
			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s\n", k, h)))
		}
	}
	_, _ = w.Write([]byte("\n"))
	_, _ = w.Write(resp.Body)

}
