package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ltto/kakaxi"
)

var root = "dao/target/angelscript.hazelight.se/"
var fixList string

func main() {
	file, err := os.ReadFile("dao/m_content.json")
	if err != nil {
		panic(err)
	}
	var links []string
	err = json.Unmarshal(file, &links)
	if err != nil {
		panic(err)
	}
	for i, link := range links {
		if i <= 1 {
			continue
		}
		fmt.Printf("\r page %d/%d++++", i, len(links))
		parts := strings.Split(link, "#")
		if len(parts) != 2 {
			continue
		}
		params := parts[1]
		paramsParts := strings.Split(params, ":")
		if len(paramsParts) != 2 {
			continue
		}
		class := paramsParts[1]
		GetBody(fmt.Sprintf("https://angelscript.hazelight.se/api/classes/C/%s-Summary.js", class))         //-Summary.js
		GetBody(fmt.Sprintf("https://angelscript.hazelight.se/api/classes/C/%s-SummaryToolTips.js", class)) //-SummaryToolTips.js
		GetBody(fmt.Sprintf("https://angelscript.hazelight.se/api/classes/C/%s-ToolTips.js", class))        //-ToolTips.js
		GetBody(fmt.Sprintf("https://angelscript.hazelight.se/api/classes/C/%s.html", class))               //.html
	}
	err = os.WriteFile("dao/fix.sh", []byte(fixList), 0777)
	if err != nil {
		panic(err)
	}
}
func GetBody(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	filePath := path.Join(root, req.URL.Path)
	//文件是否存在
	exist := kakaxi.FileExist(filePath)
	if exist {
		return
	}
	fixList += "wget " + url + " \n"

	return
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_ = os.MkdirAll(path.Dir(filePath), 0777)
	err = os.WriteFile(filePath, body, 0777)
	if err != nil {
		panic(err)
	}
}
