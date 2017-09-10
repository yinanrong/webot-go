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
	mux  map[string]http.Handler
}

func (c *controller) BadRequest(w http.ResponseWriter, err ErrorResponse) {
	c.setHead(w)
	w.WriteHeader(400)
	jerr, _ := json.Marshal(err)
	w.Write(jerr)
}
func (c *controller) StatusCode(w http.ResponseWriter, code int, err *ErrorResponse) {
	c.setHead(w)
	w.WriteHeader(code)
	jerr, _ := json.Marshal(err)
	w.Write(jerr)
}
func (c *controller) OK(w http.ResponseWriter, model interface{}) {
	c.setHead(w)
	kind := reflect.TypeOf(model).Kind()
	if kind == reflect.Struct || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Chan || kind == reflect.Map {
		json, err := json.Marshal(model)
		if err != nil {
			c.BadRequest(w, NewErrorResponse("unkown error", err.Error()))
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

type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func NewErrorResponse(code string, description string) ErrorResponse {
	return ErrorResponse{code, description}
}
