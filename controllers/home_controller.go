package controllers

import (
	"fmt"
	"net/http"
	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

type HomeController struct {
	controller
}

func (c HomeController) qr(w http.ResponseWriter, r *http.Request) {
	session, err := service.CreateSession(nil, nil)
	if err != nil {
		logs.Error(err)
		return
	}

	sessChan <- session
	sessMap[session.ID] = session
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"uuid":"%s","qr":"%s"}`, session.ID, session.Qr())
}

// r.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
// 	session, err := service.CreateSession(nil, nil)
// 	if err != nil {
// 		logs.Error(err)
// 		return
// 	}

// 	sessChan <- session
// 	sessMap[session.ID] = session
// 	w.Header().Add("Content-Type", "application/json")
// 	fmt.Fprintf(w, `{"uuid":"%s","qr":"%s"}`, session.ID, session.Qr())
// 	//w.Write(qr)
// }).Methods("GET")
// r.HandleFunc("/qr/{uuid}", func(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	if uuid, ok := vars["uuid"]; ok {
// 		if session, ok := sessMap[uuid]; ok {
// 			w.Write([]byte(strconv.Itoa(session.State)))
// 		} else {
// 			w.Write([]byte(strconv.Itoa(service.Closed)))
// 		}
// 	}
// }).Methods("GET")
// http.ListenAndServe(":5001", r)
