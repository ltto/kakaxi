package kakaxi

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

func OnTCP(conn net.Conn) error {
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Printf("http.ReadRequest err:%v\n", err)
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
