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
	"net/url"
	"strings"

	"io/ioutil"

	"fmt"

	"webot-go/service"

	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
)

// Register plugin
func Register(session *service.Session) {
	session.HandlerRegister.Add(service.MSG_TEXT, service.Handler(autoReply), "text-replier")
	if err := session.HandlerRegister.Add(service.MSG_IMG, service.Handler(autoReply), "img-replier"); err != nil {
		logs.Error(err)
	}

}
func autoReply(session *service.Session, msg *service.ReceivedMessage) {
	uri := fmt.Sprintf("http://api.qingyunke.com/api.php?key=free&appid=0&msg=%s", url.QueryEscape(msg.Content))
	resp, err := http.Get(uri)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	content, err := jc.GetString("content")
	if err != nil {
		logs.Error(err.Error())
		return
	}
	session.SendText(strings.Replace(content, "{br}", "\n", -1), session.Bot.UserName, msg.FromUserName)
}
