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

package forwarder // 以插件名命令包名

import (
	"webot-go/service" // 导入协议包

	"github.com/songtianyi/rrframework/logs" // 导入日志包
)

var (
	// 需要消息互通的群
	groups = map[string]bool{
		"jianshujiaojingdadui": true,
		"forwarder":            true, //for test
	}
)

func Register(session *service.Session) {
	session.HandlerRegister.Add(service.MSG_TEXT, service.Handler(forward), "text-forwarder")
	session.HandlerRegister.Add(service.MSG_IMG, service.Handler(forward), "img-forwarder")
}

func forward(session *service.Session, msg *service.ReceivedMessage) {
	if !msg.IsGroup {
		return
	}
	var contact *service.User
	if msg.FromUserName == session.Bot.UserName {
		contact = session.Cm.GetContactByUserName(msg.ToUserName)
	} else {
		contact = session.Cm.GetContactByUserName(msg.FromUserName)
	}
	if contact == nil {
		return
	}
	if _, ok := groups[contact.PYQuanPin]; !ok {
		return
	}
	mm, err := service.CreateMemberManagerFromGroupContact(session, contact)
	if err != nil {
		logs.Debug(err)
		return
	}
	who := mm.GetContactByUserName(msg.Who)
	if who == nil {
		who = session.Bot
	}

	for k, v := range groups {
		if !v {
			continue
		}
		c := session.Cm.GetContactByPYQuanPin(k)
		if c == nil {
			logs.Error("cannot find group contact %s", k)
			continue
		}
		if c.UserName == contact.UserName {
			// ignore
			continue
		}
		if msg.MsgType == service.MSG_TEXT {
			session.SendText("@"+who.NickName+" "+msg.Content, session.Bot.UserName, c.UserName)
		}
		if msg.MsgType == service.MSG_IMG {
			b, err := session.GetImg(msg.MsgId)
			if err == nil {
				session.SendImgFromBytes(b, msg.MsgId+".jpg", session.Bot.UserName, c.UserName)
			} else {
				logs.Error(err)
			}
		}
	}

	//mm, err := service.CreateMemberManagerFromGroupContact(contact)
}
