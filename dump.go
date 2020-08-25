package kakaxi

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func DumpRequest(r io.Reader) (header http.Header, body []byte, host, method string) {
	header = make(http.Header)
	scan := bufio.NewScanner(r)
	idx := 0
	for scan.Scan() {
		text := scan.Text()
		if idx == 0 {
			_, _ = fmt.Sscanf(text, "%s%s", &method, &host)
		} else {
			if text == "" {
				break
			}
			split := strings.Split(text, ":")
			if len(split) > 2 {
				key := strings.TrimSpace(split[0])
				v := strings.TrimSpace(split[1])
				if _, ok := header[key]; ok {
					header[key] = append(header[key], v)
				} else {
					header[key] = []string{v}
				}
			}
		}
		idx++
	}
	var length int64 = 0
	for k, hs := range header {
		if k == "Content-Length" && len(hs) > 0 {
			length, _ = strconv.ParseInt(hs[0], 10, 64)
		}
	}
	buf := bytes.NewBuffer([]byte{0})
	if length != 0 {
		_, _ = io.CopyN(buf, r, length)
	}
	return
}
