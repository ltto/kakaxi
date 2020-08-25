package kakaxi

import (
	"net"
	"net/http"
)

func OnTCP(conn net.Conn) error {
	header, body, host, method := DumpRequest(conn)
	if method == http.MethodConnect {
		return OnTLS(conn, host)
	}
	doHeader, bodyB, err := ProxyHTTP(host, method, header, body)
	if err != nil {
		return err
	}
	Writer(conn, doHeader, bodyB)
	conn.Close()
	return nil
}
