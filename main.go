package main

import (
	"time"

	"webot-go/plugins/cleaner"
	"webot-go/plugins/faceplusplus"
	"webot-go/plugins/forwarder"
	"webot-go/plugins/gifer"
	"webot-go/plugins/joker"
	"webot-go/plugins/laosj"
	"webot-go/plugins/replier"
	"webot-go/plugins/revoker"
	"webot-go/plugins/share"
	"webot-go/plugins/switcher"
	"webot-go/plugins/system"
	"webot-go/plugins/verify"
	"webot-go/plugins/youdao"
	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

func main() {
	// create session
	session, err := service.CreateSession(nil, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	// load plugins for this session
	faceplusplus.Register(session)
	replier.Register(session)
	switcher.Register(session)
	gifer.Register(session)
	cleaner.Register(session)
	laosj.Register(session)
	joker.Register(session)
	revoker.Register(session)
	forwarder.Register(session)
	system.Register(session)
	youdao.Register(session)
	verify.Register(session)
	share.Register(session)

	// enable plugin
	session.HandlerRegister.EnableByName("switcher")
	session.HandlerRegister.EnableByName("faceplusplus")
	session.HandlerRegister.EnableByName("cleaner")
	session.HandlerRegister.EnableByName("laosj")
	session.HandlerRegister.EnableByName("joker")
	session.HandlerRegister.EnableByName("system-withdraw")
	session.HandlerRegister.EnableByName("youdao")
	session.HandlerRegister.EnableByName("verify")
	session.HandlerRegister.EnableByName("share")

	// enable by type example
	session.HandlerRegister.EnableByType(service.MSG_SYS)

	//TODO:生成二维码
	//session.GenerateQR()
	for {
		if err := session.LoginAndServe(false); err != nil {
			logs.Error("session exit, %s", err)
			for i := 0; i < 3; i++ {
				logs.Info("trying re-login with cache")
				if err := session.LoginAndServe(true); err != nil {
					logs.Error("re-login error, %s", err)
				}
				time.Sleep(3 * time.Second)
			}
			if session, err = service.CreateSession(nil, session.HandlerRegister); err != nil {
				logs.Error("create new sesion failed, %s", err)
				break
			}
		} else {
			logs.Info("closed by user")
			break
		}
	}
}
