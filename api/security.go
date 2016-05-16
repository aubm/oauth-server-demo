package api

import (
	"net/http"
	"strings"

	"github.com/RangelReale/osin"
	"github.com/aubm/oauth-server-demo/security"
	"github.com/gorilla/context"
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
			if u, err := h.UsersFinder.FindByCredentials(ar.Username, ar.Password); err == nil {
				ar.UserData = u
				ar.Authorized = true
			}
		}
		h.AuthServer.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

type IdentityAdapter struct {
	AccessFinder interface {
		FindAccess(code string) (*osin.AccessData, error)
	} `inject:""`
}

func (ia *IdentityAdapter) Adapt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var u *security.User
		accessToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if accessToken != "" {
			if accessData, err := ia.AccessFinder.FindAccess(accessToken); err == nil {
				if userData, ok := accessData.UserData.(*security.User); ok {
					u = userData
				}
			}
		}
		if u == nil {
			httpError(w, 403, "invalid_token", "Invalid access token")
			return
		}
		context.Set(r, "user", u)
		next.ServeHTTP(w, r)
	})
}
