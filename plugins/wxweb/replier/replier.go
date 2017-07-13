/*
Copyright 2017 wechat-go Authors. All Rights Reserved.
MIT License

Copyright (c) 2017

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package replier

import (
	"net/http"

	"io/ioutil"

	"fmt"

	"encoding/json"

	"github.com/songtianyi/rrframework/logs"
	"github.com/yinanrong/wechat-go/wxweb"
)

// register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(autoReply), "text-replier")
	if err := session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(autoReply), "img-replier"); err != nil {
		logs.Error(err)
	}

}
func autoReply(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	url := fmt.Sprintf("http://api.qingyunke.com/api.php?key=free&appid=0&msg=%s", msg.Content)
	resp, err := http.Get(url)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	var r replay
	err = json.Unmarshal(body, &r)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	session.SendText(r.Content, session.Bot.UserName, msg.FromUserName)
	// if !msg.IsGroup {
	// 	session.SendText("暂时不在，稍后回复", session.Bot.UserName, msg.FromUserName)
	// }
}

type replay struct {
	Result  int    `json:"result"`
	Content string `json:"content"`
}
