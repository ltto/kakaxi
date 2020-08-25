package kakaxi

import (
	"net"
	"net/http"
)

func OnTCP(conn net.Conn) error {
	header, host, method := DumpRequest(conn)
	if method == http.MethodConnect {
		return OnTLS(conn, host)
	}
	doHeader, bodyB, err := ProxyHTTP(host, method, header)
	if err != nil {
		return err
	}
	Writer(conn, doHeader, bodyB)
	conn.Close()
	return nil
}
