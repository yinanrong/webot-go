package main

import (
	"net/http"
	"webot-go/plugins/replier"
	"webot-go/plugins/switcher"
	"webot-go/plugins/system"
	"webot-go/plugins/verify"
	"webot-go/plugins/youdao"
	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

func main() {

	sesChan := make(chan (*service.Session), 100)

	go apiService(sesChan)

	for s := range sesChan {
		go backService(s)
	}
	select {}
}
func apiService(sesChan chan<- (*service.Session)) {
	http.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
		session, err := service.CreateSession(nil, nil)
		if err != nil {
			logs.Error(err)
			return
		}
		qr, err := session.GenerateQR()
		if err != nil {
			logs.Error(err)
			return
		}
		sesChan <- session
		w.Write(qr)
	})
	http.ListenAndServe(":1989", nil)
}

func backService(session *service.Session) {
	switcher.Register(session)
	replier.Register(session)
	system.Register(session)
	youdao.Register(session)
	verify.Register(session)

	session.HandlerRegister.EnableByName("switcher")
	session.HandlerRegister.EnableByName("youdao")
	session.HandlerRegister.EnableByName("system-sys")
	session.HandlerRegister.EnableByName("system-withdraw")
	session.HandlerRegister.EnableByName("verify")

	if session.LoginAndServe(false) != nil {
		logs.Info("closed by user")
	}
}
