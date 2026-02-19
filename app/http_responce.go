package app

import (
	"encoding/json"
	"net/http"
	"text/template"
)

type ErrorJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func respondWithStructJSON(w http.ResponseWriter, code int, data interface{}) {
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(data)
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	resp := ErrorJson{Code: statusCode, Message: err.Error()}
	response, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}

func respondWithTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if data != nil {
		_ = tmpl.Execute(w, data)
	}
}
