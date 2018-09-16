package serverAdapter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"

	"../picConverter"
	"../weixinAdapter"
	"golang.org/x/net/websocket"
)

type AccountInfo struct {
	LoginPic string
	Uuid     string
}

func qrhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端

		qrPic, uuid := weixinAdapter.GetQrPic()
		account := AccountInfo{LoginPic: picConverter.Convert2Base64(qrPic), Uuid: uuid}
		jsonStr, _ := json.Marshal(account)
		io.WriteString(w, string(jsonStr))
	}
}

func newpagehandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端

		m, _ := url.ParseQuery(r.URL.RawQuery)
		newpageModel := weixinAdapter.NewpageModel{Ticket: m["ticket"][0], Uuid: m["uuid"][0], Scan: m["scan"][0]}
		uri := weixinAdapter.NewPage(newpageModel)

		io.WriteString(w, string(uri))
	}
}

func tryloghandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端

		m, _ := url.ParseQuery(r.URL.RawQuery)

		uri := weixinAdapter.TryLogin(weixinAdapter.TryLoginModel{Uuid: m["uuid"][0]})

		io.WriteString(w, uri)
	}
}

func web(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法

	fmt.Println("method", r.Method)

	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端

		t, _ := template.ParseFiles("websocket.html")

		t.Execute(w, nil)

	} else {

		//否则走打印输出post接受的参数username和password

		fmt.Println(r.PostFormValue("username"))

		fmt.Println(r.PostFormValue("password"))

	}

}

func StartServer() {
	http.Handle("/websocket", websocket.Handler(echo))
	http.HandleFunc("/web", web)
	http.HandleFunc("/QrCode", qrhandler)
	http.HandleFunc("/TryLogin", tryloghandler)
	http.HandleFunc("/NewPage", newpagehandler)

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func echo(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed:", err)
			break
		}

		fmt.Println("received from client:" + reply)
		msg := "received:" + reply
		fmt.Println("send to client:" + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("send failed:", err)
			break
		}

	}
}
