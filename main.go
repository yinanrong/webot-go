package main

import (
	"net/http"
	"time"
	"webot-go/controllers"
	"webot-go/plugins/replier"
	"webot-go/plugins/verify"
	"webot-go/service"

	"github.com/astaxie/beego/logs"
)

func main() {
	sessMap := make(map[string]*service.Session)
	sessChan := make(chan (*service.Session), 100)
	service.InitSessVector(sessMap, sessChan)
	go http.ListenAndServe(":5001", controllers.NewHomeController())
	go func() {
		for range time.Tick(time.Hour) {
			for uuid, session := range sessMap {
				if session.State == service.Closed {
					delete(sessMap, uuid)
				}
			}
		}
	}()
	go func() {
		for session := range sessChan {
			go backService(session, sessMap)
		}
		close(sessChan)
	}()
	select {}
}

func backService(session *service.Session, sessMap map[string]*service.Session) {
	replier.Register(session)
	verify.Register(session)
	session.HandlerRegister.EnableByName("text-replier")
	session.HandlerRegister.EnableByName("verify")
	if err := session.LoginAndServe(); err != nil {
		logs.Info("session closed due to :%s", err.Error())
		time.Sleep(5 * time.Second)
		session.State = service.Closed
	}
}
