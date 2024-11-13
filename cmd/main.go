package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/ltto/kakaxi"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	port := *flag.Int("p", 8081, "端口")
	host := *flag.String("h", "", "绑定主机地址")
	flag.Parse()
	sprintf := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("监听", sprintf)
	listen, err := net.Listen("tcp", sprintf)
	if err != nil {
		panic(err)
	}
	for true {
		accept, err := listen.Accept()
		if err != nil {
			continue
		}
		go func() {
			kakaxi.OnTCP(accept)
		}()
	}
}
