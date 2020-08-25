package kakaxi

import "net/http"

var OnRequest = func(header http.Header, host string, body []byte) {}
var OnResponse = func(header http.Header, host string, body []byte) {}
