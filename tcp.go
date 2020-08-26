package kakaxi

import (
	"bufio"
	"net"
	"net/http"
)

func OnTCP(conn net.Conn) error {
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}
	if request.Method == http.MethodConnect {
		return OnTLS(conn)
	}

	doHeader, resp, bodyB, err := ProxyHTTP(*request)
	if err != nil {
		return err
	}
	Writer(conn, doHeader, resp, bodyB)
	conn.Close()
	return nil
}
