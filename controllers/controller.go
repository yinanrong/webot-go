package controllers

import "net/http"

import "fmt"
import "encoding/json"
import "reflect"

const (
	views = "../views/"
)

type controller struct {
	Name string
	mux  map[string]func(w http.ResponseWriter, r *http.Request)
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := c.mux[r.URL.Path]; ok {
		h(w, r)
		return
	}
	http.ServeFile(w, r, "../views/notfound.html")
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

func (*controller) View(view string, w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, views+view)
}
