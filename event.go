package kakaxi

import "net/http"

type OnRequest = func(header http.Header, host string, body []byte)
type OnResponse = func()
