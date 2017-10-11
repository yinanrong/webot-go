package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"webot-go/service"

	"github.com/gorilla/mux"
	"github.com/songtianyi/rrframework/logs"
)

type HomeController struct {
	controller
}

func NewHomeController() *HomeController {
	h := new(HomeController)
	if h.mux == nil {
		h.mux = make(map[string]func(w http.ResponseWriter, r *http.Request))
	}
	h.mux["/index"] = h.index
	h.mux["/qr"] = h.qr
	h.mux["/state"] = h.state
	return h
}
func (c *HomeController) index(w http.ResponseWriter, r *http.Request) {
	c.View("index.html", w, r)
}
func (c *HomeController) qr(w http.ResponseWriter, r *http.Request) {
	session, err := service.CreateSession(nil, nil)
	if err != nil {
		logs.Error(err)
		c.BadRequest(w, err)
		return
	}
	session.EnQueue()
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"uuid":"%s","qr":"%s"}`, session.ID, session.Qr())
}
func (c *HomeController) state(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if uuid, ok := vars["uuid"]; ok {
		if session, ok := service.SessionPop(uuid); ok {
			w.Write([]byte(strconv.Itoa(session.State)))
		} else {
			w.Write([]byte(strconv.Itoa(service.Closed)))
		}
	}
}
