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

	reps, err := ProxyHTTP(CopyRequest(request))
	if err != nil {
		return err
	}
	Writer(conn, reps)
	_ = conn.Close()
	return nil
}
