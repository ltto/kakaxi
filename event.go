package kakaxi

import "net/http"

var OnRequest = func(r http.Request) {}
var OnResponse = func(r http.Request, header http.Header, body []byte) {}
