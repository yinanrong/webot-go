package controllers

import "net/http"

import "fmt"
import "encoding/json"
import "reflect"

type controller struct {
	mux map[string]func(w http.ResponseWriter, r *http.Request)
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := c.mux[r.URL.Path]; ok {
		h(w, r)
		return
	}
	c.NotFound(w, r)
}
func (c *controller) BadRequest(w http.ResponseWriter, msg interface{}) {
	c.setHead(w)
	w.WriteHeader(400)
	json, _ := json.Marshal(msg)
	w.Write(json)
}
func (c *controller) StatusCode(w http.ResponseWriter, code int, msg interface{}) {
	c.setHead(w)
	w.WriteHeader(code)
	json, _ := json.Marshal(msg)
	w.Write(json)
}
func (c *controller) OK(w http.ResponseWriter, model interface{}) {
	c.setHead(w)
	kind := reflect.TypeOf(model).Kind()
	if kind == reflect.Struct || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Chan || kind == reflect.Map {
		json, err := json.Marshal(model)
		if err != nil {
			c.BadRequest(w, err.Error())
		} else {
			w.Write(json)
		}
	} else {
		fmt.Fprint(w, model)
	}
}
func (*controller) setHead(w http.ResponseWriter) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
}

func (*controller) View(content string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

func (*controller) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("请求的地址" + r.URL.Path + "不存在"))
}
