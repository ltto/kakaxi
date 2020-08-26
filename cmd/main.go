package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/ltto/kakaxi"
)

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
			if err = kakaxi.OnTCP(accept); err != nil {
				fmt.Println("wocao ", err)
			}
		}()
	}
}
