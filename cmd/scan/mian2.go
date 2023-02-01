package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	//
	fmt.Println("Hello, playground")
	Dir("/Users/ltt/ltto/kakaxi/dao/target/cocafe.co")
}

func Dir(path string) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			Dir(path + "/" + dir.Name())
		} else {
			if dir.Size() == 0 {
				//fmt.Println(path + "/" + dir.Name())
			}
			if strings.HasSuffix(dir.Name(), ".meta.json") {
				err := os.Remove(path + "/" + dir.Name())
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
