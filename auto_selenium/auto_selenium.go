package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ltto/selenium"
	selenium2 "github.com/tebeka/selenium"

	"github.com/ltto/kakaxi"
)

var proxyHttp = fmt.Sprintf("127.0.0.1:%d", 8088)

func runProxy(ch chan<- struct{}) {
	listen, err := net.Listen("tcp", proxyHttp)
	if err != nil {
		panic(err)
	}
	ch <- struct{}{}

	for true {
		accept, err := listen.Accept()
		if err != nil {
			continue
		}
		go func() {
			_ = kakaxi.OnTCP(accept)
		}()
	}
}
func init() {
	ch := make(chan struct{})
	go runProxy(ch)
	<-ch
	fmt.Printf("proxy server start at %s\n", proxyHttp)
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	// 设置HTTP代理
	caps := selenium.NewCapabilities()
	driver, err := selenium.NewChromeDriverCap(*caps, selenium.ChromeCapabilities{Args: []string{
		"--proxy-server=http://" + proxyHttp,
		"--root-ca-cert=D:\\workspace\\sync_work\\GolandProjects\\kakaxi\\kakaxi-ca-cert.pem",
		"--ignore-certificate-errors",           // 忽略证书错误
		"--ignore-ssl-errors",                   // 忽略 SSL 错误
		"--ignore-certificate-errors-spki-list", // 忽略特定证书错误
		"--allow-insecure-localhost",            // 允许不安全的本地连接
		"--disable-web-security",                // 禁用 Web 安全性检查
		//"--headless", "--no-sandbox", "--disable-gpu", // 无头模式
	}})
	if err != nil {
		panic(err)
	}
	var links = []string{
		"https://angelscript.hazelight.se/",
		"https://angelscript.hazelight.se/getting-started/installation/",
		"https://angelscript.hazelight.se/getting-started/introduction/",
		"https://angelscript.hazelight.se/scripting/functions-and-events/",
		"https://angelscript.hazelight.se/scripting/properties-and-accessors/",
		"https://angelscript.hazelight.se/scripting/actors-components/",
		"https://angelscript.hazelight.se/scripting/function-libraries/",
		"https://angelscript.hazelight.se/scripting/fname-literals/",
		"https://angelscript.hazelight.se/scripting/format-strings/",
		"https://angelscript.hazelight.se/scripting/structs-refs/",
		"https://angelscript.hazelight.se/scripting/networking-features/",
		"https://angelscript.hazelight.se/scripting/delegates/",
		"https://angelscript.hazelight.se/scripting/mixin-methods/",
		"https://angelscript.hazelight.se/scripting/gameplaytags/",
		"https://angelscript.hazelight.se/scripting/editor-script/",
		"https://angelscript.hazelight.se/scripting/subsystems/",
		"https://angelscript.hazelight.se/scripting/script-tests/",
		"https://angelscript.hazelight.se/scripting/cpp-differences/",
		"https://angelscript.hazelight.se/cpp-bindings/automatic-bindings/",
		"https://angelscript.hazelight.se/cpp-bindings/mixin-libraries/",
		"https://angelscript.hazelight.se/cpp-bindings/precompiled-data/",
		"https://angelscript.hazelight.se/project/development-status/",
		"https://angelscript.hazelight.se/project/resources/",
		"https://angelscript.hazelight.se/project/license/",
	}
	//file, err := os.ReadFile("dao/m_content.json")
	//if err != nil {
	//	panic(err)
	//}
	//err = json.Unmarshal(file, &links)
	//if err != nil {
	//	panic(err)
	//}
	length := len(links)
	for i := 0; i < len(links); i++ {
		link := links[i] //正序
		fmt.Printf("\r page %d/%d++++", i, length)
		driver.Get(link)
		//time.Sleep(time.Second)
	}
	driver.Quit()
}

// WaitForIframe 等待 iframe 加载并切换到 iframe
func WaitForIframe(driver *selenium.WebDriver, by, value string, timeout time.Duration) error {
	condition := func(wd selenium2.WebDriver) (bool, error) {
		// 尝试找到 iframe
		frames, err := wd.FindElements(by, value)
		if err != nil || len(frames) == 0 {
			return false, nil
		}

		// 尝试切换到 iframe
		if err := wd.SwitchFrame(frames[0]); err != nil {
			return false, nil
		}

		// 切换回主文档
		wd.SwitchFrame(nil)
		return true, nil
	}

	// 使用 selenium 提供的 WaitWithTimeout
	driver.WaitWithTimeout(condition, timeout)

	// 最终切换到 iframe
	frames := driver.FindElements(by, value)
	driver.SwitchFrame(frames[0])

	return nil
}
