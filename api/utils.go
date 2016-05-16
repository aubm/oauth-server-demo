package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/context"
)

const SERVER_ERR = "server_error"
const SERVER_ERR_DESC = "An error occured, please try again later"

func readJSON(r io.Reader, w http.ResponseWriter, target interface{}) error {
	var err error
	b, err := ioutil.ReadAll(r)
	if err == nil {
		err = json.Unmarshal(b, target)
	}
	if err != nil {
		http.Error(w, "Invalid JSON input", 400)
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	b, err := json.Marshal(data)
	if err != nil {
		httpError(w, 500, SERVER_ERR, SERVER_ERR_DESC)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func httpError(w http.ResponseWriter, code int, msg, description string) {
	writeJSON(w, struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{msg, description}, code)
}

var validEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func validateEmailFormat(email string) bool {
	return validEmail.MatchString(email)
}

type ClearContextAdapter struct{}

func (a *ClearContextAdapter) Adapt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		next.ServeHTTP(w, r)
	})
}
