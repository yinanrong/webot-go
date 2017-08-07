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

package switcher

import (
	"strings"

	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

// Register plugin
func Register(session *service.Session) {
	session.HandlerRegister.Add(service.MSG_TEXT, service.Handler(switcher), "switcher")
}

func switcher(session *service.Session, msg *service.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)

	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}
	if contact.UserName != session.Bot.UserName {
		session.SendText("hehe, too much you think. only @"+session.Bot.NickName+" can use this function", session.Bot.UserName, service.RealTargetUserName(session, msg))
	}

	if strings.ToLower(msg.Content) == "dump" {
		session.SendText(session.HandlerRegister.Dump(), session.Bot.UserName, service.RealTargetUserName(session, msg))
		return
	}

	if !strings.Contains(strings.ToLower(msg.Content), "enable") &&
		!strings.Contains(strings.ToLower(msg.Content), "disable") {
		return
	}

	ss := strings.Split(msg.Content, " ")
	if len(ss) < 2 {
		return
	}
	if strings.ToLower(ss[1]) == "switcher" {
		session.SendText("hehe,too much you think", session.Bot.UserName, service.RealTargetUserName(session, msg))
		return
	}

	var (
		err error
	)
	if strings.ToLower(ss[0]) == "enable" {
		if err = session.HandlerRegister.EnableByName(ss[1]); err == nil {
			session.SendText(msg.Content+" [DONE]", session.Bot.UserName, service.RealTargetUserName(session, msg))
		} else {
			session.SendText(err.Error(), session.Bot.UserName, service.RealTargetUserName(session, msg))
		}
	} else if strings.ToLower(ss[0]) == "disable" {
		if err = session.HandlerRegister.DisableByName(ss[1]); err == nil {
			session.SendText(msg.Content+" [DONE]", session.Bot.UserName, service.RealTargetUserName(session, msg))
		} else {
			session.SendText(err.Error(), session.Bot.UserName, service.RealTargetUserName(session, msg))
		}
	}
}
