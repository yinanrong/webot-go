package forwarder

import (
	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

var (
	// 需要消息互通的群
	groups = map[string]bool{
		"a": true,
		"b": true,
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
