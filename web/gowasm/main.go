package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"syscall/js"
	"github.com/JackZxj/wasmgo-demo/web/gowasm/goshell"
)

var (
	shell *goshell.SecureShell

	document = js.Global().Get("document")

	hostElm     = document.Call("getElementById", "host")
	portElm     = document.Call("getElementById", "port")
	userElm     = document.Call("getElementById", "user")
	passwordElm = document.Call("getElementById", "password")
	errElm      = document.Call("getElementById", "err")

	connectBtn    = js.Global().Get("connect")
	disconnectBtn = js.Global().Get("disconnect")
	execBtn       = js.Global().Get("exec")

	resultElm = document.Call("getElementById", "result")
)

func connectJs(this js.Value, args []js.Value) interface{} {
	var (
		err error
		buf bytes.Buffer

		host     string
		port     int
		user     string
		password string
	)

	connectBtn.Set("disabled", js.ValueOf(true))

	hostValue := hostElm.Get("value")
	portValue := portElm.Get("value")
	userValue := userElm.Get("value")
	passwordValue := passwordElm.Get("value")

	if host = hostValue.String(); host == "" {
		errElm.Set("innerHTML", js.ValueOf("IP can not be empty!"))
		connectBtn.Set("disabled", js.ValueOf(false))
		return nil
	}

	if port, err = strconv.Atoi(portValue.String()); err != nil {
		errElm.Set("innerHTML", js.ValueOf(fmt.Sprintf("can not convert port to int: %v", err)))
		connectBtn.Set("disabled", js.ValueOf(false))
		return nil
	} else if port == 0 {
		port = 22
		portElm.Set("value", js.ValueOf(port))
	}

	if user = userValue.String(); user == "" {
		user = "root"
		userElm.Set("value", js.ValueOf("root"))
	}

	if password = passwordValue.String(); password == "" {
		errElm.Set("innerHTML", js.ValueOf("password can not be empty!"))
		connectBtn.Set("disabled", js.ValueOf(false))
		return nil
	}

	shell, err = goshell.NewSecureShell(&buf, host, user, password, port)
	if err != nil {
		errElm.Set("innerHTML", js.ValueOf(fmt.Sprintf("err: %v", err)))
		connectBtn.Set("disabled", js.ValueOf(false))
		return nil
	}
	execBtn.Set("disabled", js.ValueOf(false))
	disconnectBtn.Set("disabled", js.ValueOf(false))
	innerHtml := resultElm.Get("innerHTML")
	resultElm.Set("innerHTML", js.ValueOf(fmt.Sprintf("%s\nLogin %q success!\n%s", innerHtml, host, buf.String())))

	return nil
}

func disconnectJs(this js.Value, args []js.Value) interface{} {
	shell = nil
	execBtn.Set("disabled", js.ValueOf(true))
	disconnectBtn.Set("disabled", js.ValueOf(true))
	connectBtn.Set("disabled", js.ValueOf(false))
	return nil
}

func main() {
	fmt.Println("Hello, WASM!") // 会作为输出到控制台
	res, err := http.Get("index.html") // 只能获取同一 server 下的页面，否则会跨域
	if err != nil {
		fmt.Println("get http", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println("read err", err)
	}
	fmt.Println("get body:", string(body))

	if _, err = net.Dial("tcp", "10.110.26.178:8080"); err != nil {
		fmt.Println("dial err:", err)
	}

	done := make(chan int, 0)
	connectBtn.Call("addEventListener", "click", js.FuncOf(connectJs))
	disconnectBtn.Call("addEventListener", "click", js.FuncOf(disconnectJs))
	<-done
}
