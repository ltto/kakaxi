package main

import (
	"fmt"
	"net"

	"github.com/ltto/kakaxi"
)

func main() {
	listen, err := net.Listen("tcp", ":8081")
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
