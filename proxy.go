package kakaxi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func ProxyHTTP(req Request) (resp Response, err error) {
	URL := req.URL
	savePath := URL.Host + URL.Path
	if path.Ext(savePath) == "" {
		savePath = savePath + "/index.html"
	}
	savePath = path.Join("dao/target", savePath)
	exist := FileExist(savePath)
	if exist {
		if resp, err = FileToHttpResponse(savePath); err != nil {
			return
		}
	} else {
		if resp, err = doProxyHTTP(req); err != nil {
			return
		}
	}
	go OnResponse(req, resp, !exist, savePath, URL.Host, URL.Path)
	return
}
func doProxyHTTP(r Request) (resp Response, err error) {
	r.Request.RequestURI = ""
	var do *http.Response
	do, err = http.DefaultClient.Do(r.Request)
	if err != nil {
		log.Println("http.DefaultClient.Do", err)
		return
	}
	resp = CopyResponse(do)
	return
}

func Writer(w io.Writer, resp Response) {
	_, _ = w.Write([]byte(fmt.Sprintf("%s %s\n", resp.Proto, resp.Status)))
	for k, hs := range resp.Header {
		for _, h := range hs {
			_, _ = w.Write([]byte(fmt.Sprintf("%s: %s\n", k, h)))
		}
	}
	_, _ = w.Write([]byte("\n"))
	_, _ = w.Write(resp.Body)
}

func FileToHttpResponse(filePath string) (Response, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("os.Open(%s) err:%v\n", filePath, err)
		return Response{}, err
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("io.ReadAll(%s) err:%v\n", filePath, err)
		return Response{}, err
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("file.Stat(%s) err:%v\n", filePath, err)
		return Response{}, err
	}

	// 获取文件类型
	contentType := http.DetectContentType(content)

	// 构造 Response
	response := Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		Body:       content,
		Header: http.Header{
			"Content-Type":   []string{contentType},
			"Content-Length": []string{string(fileInfo.Size())},
			"Last-Modified":  []string{fileInfo.ModTime().Format(time.RFC1123)},
			"Accept-Ranges":  []string{"bytes"},
		},
	}

	return response, nil
}
