package main

import (
	"encoding/json"
	"net/http"
)

type J map[string]interface{}

type Web struct {
	w http.ResponseWriter
}

func (p *Web) Json(code int, data interface{}) (int, string) {
	b, err := json.Marshal(data)
	chk(err)
	p.w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return code, string(b)
}

func (p *Web) Code(code int) (int, string) {
	return code, http.StatusText(code)
}