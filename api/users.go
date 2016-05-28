package api

import (
	"net/http"

	"github.com/aubm/oauth-server-demo/security"
	"github.com/gorilla/context"
)

type UsersHandlers struct {
	Manager interface {
		Save(user security.User) error
		FindByEmail(email string) (*security.User, error)
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
	if !validateEmailFormat(u.Email) {
		httpError(w, 400, "email_format_error", "Invalid email format")
		return
	}
	if _, err := h.Manager.FindByEmail(u.Email); err == nil {
		httpError(w, 400, "email_unicity_error", "This email is already used")
		return
	} else {
		if _, ok := err.(security.NoUserFoundErr); !ok {
			httpError(w, 500, SERVER_ERR, SERVER_ERR_DESC)
			return
		}
	}

	if err := h.Manager.Save(u); err != nil {
		httpError(w, 500, SERVER_ERR, SERVER_ERR_DESC)
		return
	}
	w.WriteHeader(201)
}

func (h *UsersHandlers) Me(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*security.User)
	writeJSON(w, map[string]string{"id": user.Id, "email": user.Email}, 200)
}
