package api

import (
	"net/http"

	"github.com/aubm/oauth-server-demo/security"
)

type UsersHandlers struct {
	Manager interface {
		Save(user security.User) error
	} `inject:""`
}

func (h *UsersHandlers) Create(w http.ResponseWriter, r *http.Request) {
	u := security.User{}
	if err := readJSON(r.Body, w, &u); err != nil {
		return
	}
	if u.Email == "" || u.Password == "" {
		httpError(w, 400, "form_error", "A user must have an email and a password")
		return
	}
	if err := h.Manager.Save(u); err != nil {
		httpError(w, 500, SERVER_ERR, SERVER_ERR_DESC)
		return
	}
	w.WriteHeader(201)
}
