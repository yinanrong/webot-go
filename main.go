package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webot-go/plugins/replier"
	"webot-go/plugins/verify"
	"webot-go/service"

	"github.com/gorilla/mux"
	"github.com/songtianyi/rrframework/logs"
)

func main() {
	sessMap := make(map[string]*service.Session)
	sessChan := make(chan (*service.Session), 100)
	go apiService(sessChan, sessMap)

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
func apiService(sessChan chan<- *service.Session, sessMap map[string]*service.Session) {
	r := mux.NewRouter()
	r.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
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
		sessChan <- session
		sessMap[session.ID] = session
		w.Header().Add("Content-Type", "application/json")
		base64Qr := base64.StdEncoding.EncodeToString(qr)
		w.Write([]byte(fmt.Sprintf(`{"uuid":"%s","qr":"%s"}`, session.ID, base64Qr)))
		//w.Write(qr)
	}).Methods("GET")
	r.HandleFunc("/qr/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if uuid, ok := vars["uuid"]; ok {
			if session, ok := sessMap[uuid]; ok {
				w.Write([]byte(strconv.Itoa(session.State)))
			} else {
				w.Write([]byte(strconv.Itoa(service.Closed)))
			}
		}
	}).Methods("GET")
	http.ListenAndServe(":5001", r)
}

func backService(session *service.Session, sessMap map[string]*service.Session) {
	replier.Register(session)
	//youdao.Register(session)
	verify.Register(session)

	// session.HandlerRegister.EnableByName("youdao")
	session.HandlerRegister.EnableByName("text-replier")
	session.HandlerRegister.EnableByName("verify")
	if err := session.LoginAndServe(); err != nil {
		logs.Info("session closed due to :%s", err.Error())
		time.Sleep(5 * time.Second)
		session.State = service.Closed
	}
}
