package kakaxi

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DumpRequest(r io.Reader) (header http.Header, host, method string) {
	header = make(http.Header)
	scan := bufio.NewScanner(r)
	idx := 0
	for scan.Scan() {
		text := scan.Text()
		if idx == 0 {
			fmt.Sscanf(text, "%s%s", &method, &host)
		} else {
			if text == "" {
				return
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
	return
}
