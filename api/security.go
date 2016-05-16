package api

import (
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/aubm/oauth-server-demo/security"
)

type SecurityHandlers struct {
	AuthServer interface {
		NewResponse() *osin.Response
		HandleAccessRequest(w *osin.Response, r *http.Request) *osin.AccessRequest
		FinishAccessRequest(w *osin.Response, r *http.Request, ar *osin.AccessRequest)
	} `inject:""`
	UsersFinder interface {
		FindByCredentials(email, clearPassword string) (*security.User, error)
	} `inject:""`
}

func (h *SecurityHandlers) Token(w http.ResponseWriter, r *http.Request) {
	resp := h.AuthServer.NewResponse()
	defer resp.Close()

	if ar := h.AuthServer.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.PASSWORD:
			if _, err := h.UsersFinder.FindByCredentials(ar.Username, ar.Password); err == nil {
				ar.Authorized = true
			}
		}
		h.AuthServer.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}
