package weixinAdapter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../fileAdapter"
)

type NewpageModel struct {
	Ticket string
	Uuid   string
	Scan   string
}

type TryLoginModel struct {
	Uuid string
}

type NewpageResult struct {
	ret         string
	message     string
	skey        string
	wxsid       string
	wxuin       string
	pass_ticket string
	isgrayscale string
	deviceId    string
}

func GetQrPic() ([]byte, string) {
	t := time.Now()
	//调用微信页面
	resp, err := http.Get("https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + string(t.Unix()))
	if err != nil {
		fmt.Println(err)
	}
	//获取登录图片
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	pattern := `window.QRLogin.uuid = ".*";`
	reg, err := regexp.Compile(pattern)
	aimSlice := reg.FindString(string(body))
	uuid := aimSlice[23:35]

	resp, err = http.Get(`https://login.weixin.qq.com/qrcode/` + uuid)

	body, err = ioutil.ReadAll(resp.Body)

	return body, uuid
}

func NewPage(model NewpageModel) []byte {
	requestUrl := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage" +
		"?ticket=" + model.Ticket +
		"&uuid=" + model.Uuid +
		"&scan=" + model.Scan +
		"&lang=zh_CN&fun=new"
	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
	return body
}

func TryLogin(model TryLoginModel) string {
	//https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login
	//https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=wfhCwDHksg==&tip=0&r=594480321&_=1537003786352
	resp, err := http.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=0&uuid=" + model.Uuid + "&_=" + string(time.Now().Unix()))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
	pattern := `window.code=.*;`
	reg, err := regexp.Compile(pattern)
	aimSlice := reg.FindString(string(body))
	resultCode := aimSlice[12:15]

	fmt.Println(resultCode)

	if resultCode != "200" {
		return ""
	} else {
		pattern = `window.redirect_uri=.*;`
		reg, err = regexp.Compile(pattern)
		newpageUrl := reg.FindString(string(body))
		newpageUrl = strings.Replace(newpageUrl, `window.redirect_uri="`, ``, -1)
		newpageUrl = strings.Replace(newpageUrl, `";`, ``, -1)
		u, _ := url.Parse(newpageUrl)
		values, _ := url.ParseQuery(u.RawQuery)

		fmt.Println(values)
		//window.redirect_uri="https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AbZwT1Nfy85lVUPrashsYTkQ@qrticket_0&uuid=YdbzRbsbOw==&lang=zh_CN&scan=1537004878";

		newpageResult := NewPage(NewpageModel{Ticket: values["ticket"][0], Uuid: values["uuid"][0], Scan: values["scan"][0]})

		pattern = `<pass_ticket>.*</pass_ticket>`
		reg, err = regexp.Compile(pattern)
		passTicket := reg.FindString(string(newpageResult))
		passTicket = strings.Replace(passTicket, `<pass_ticket>`, ``, -1)
		passTicket = strings.Replace(passTicket, `</pass_ticket>`, ``, -1)

		pattern = `<skey>.*</skey>`
		reg, err = regexp.Compile(pattern)
		sKey := reg.FindString(string(newpageResult))
		sKey = strings.Replace(sKey, `<skey>`, ``, -1)
		sKey = strings.Replace(sKey, `</skey>`, ``, -1)

		pattern = `<wxsid>.*</wxsid>`
		reg, err = regexp.Compile(pattern)
		sId := reg.FindString(string(newpageResult))
		sId = strings.Replace(sId, `<wxsid>`, ``, -1)
		sId = strings.Replace(sId, `</wxsid>`, ``, -1)

		pattern = `<wxuin>.*</wxuin>`
		reg, err = regexp.Compile(pattern)
		uIn := reg.FindString(string(newpageResult))
		uIn = strings.Replace(uIn, `<wxuin>`, ``, -1)
		uIn = strings.Replace(uIn, `</wxuin>`, ``, -1)

		fmt.Println(passTicket)
		fmt.Println(sKey)

		initResult := webwxinit(NewpageResult{pass_ticket: passTicket, skey: sKey, wxsid: sId, wxuin: uIn, deviceId: values["scan"][0]})

		return string(initResult)

	}
}

func webwxinit(model NewpageResult) []byte {

	// 、https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit?r=-1766680381&lang=zh_CN&pass_ticket=aaZXuDcg6yE6bOqFMBQkc9XcnV54VCX1PkiYRkqvytvjsp9ssy6TkmOt4t2%252BYoVR
	requestUrl := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit" +
		"?pass_ticket=" + model.pass_ticket +
		//"&skey=" + model.skey +
		//"&lang=zh_CN" +
		"&r=" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:10]

	fmt.Println(requestUrl)

	postParam := "{\"BaseRequest\": { " +
		"\"Uin\":\"" + model.wxuin + "\"," +
		"\"Sid\":\"" + model.wxsid + "\"," +
		"\"Skey\":\"" + model.skey + "\"," +
		"\"DeviceID\":\"e161526049243567\"" +
		"}}"

	fmt.Println(postParam)

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer([]byte(postParam)))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fileSucc, fErr := fileAdapter.WriteFile("successResponse.txt", body)
	if !fileSucc {
		fmt.Println(fErr.Error())
	}

	return body
}
