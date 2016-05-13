package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RangelReale/osin"
)

type authServer interface {
	NewResponse() *osin.Response
	HandleAccessRequest(w *osin.Response, r *http.Request) *osin.AccessRequest
	FinishAccessRequest(w *osin.Response, r *http.Request, ar *osin.AccessRequest)
}

type accessLoader interface {
	LoadAccess(code string) (*osin.AccessData, error)
}

type SecurityHandlers struct {
	AuthServer   authServer   `inject:""`
	AccessLoader accessLoader `inject:""`
}

func (h *SecurityHandlers) Token(w http.ResponseWriter, r *http.Request) {
	resp := h.AuthServer.NewResponse()
	defer resp.Close()

	if ar := h.AuthServer.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if ar.Username == "test" && ar.Password == "test" {
				ar.Authorized = true
			}
		}
		h.AuthServer.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func (h *SecurityHandlers) Me(w http.ResponseWriter, r *http.Request) {
	token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if token == "" {
		http.Error(w, "Invalid access token", 403)
		return
	}
	access, err := h.AccessLoader.LoadAccess(token)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	data, _ := json.Marshal(access)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
