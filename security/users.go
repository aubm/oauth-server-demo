package security

import (
	"database/sql"
	"errors"

	"github.com/aubm/oauth-server-demo/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Email    string
	Password string
}

type UsersManager struct {
	DB     *sql.DB     `inject:""`
	Config *config.App `inject:""`
}

func (m *UsersManager) Save(u User) error {
	if u.Id == 0 {
		var err error
		u.Password, err = m.encrypt(u.Password)
		if err != nil {
			return err
		}
	}
	stmt, _ := m.DB.Prepare(`INSERT INTO users
	(email, password) VALUES (?, ?)`)
	_, err := stmt.Exec(u.Email, u.Password)
	return err
}

func (m *UsersManager) FindByCredentials(email, clearPassword string) (*User, error) {
	rows, err := m.DB.Query(`SELECT id, email, password
	FROM users
	WHERE email = ?
	LIMIT 1`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No user found")
	}

	u := new(User)
	if err := rows.Scan(&u.Id, &u.Email, &u.Password); err != nil {
		return nil, err
	}

	if err := m.compareHashAndPassword(u.Password, clearPassword); err != nil {
		return nil, err
	}

	return u, nil
}

func (m *UsersManager) encrypt(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(m.Config.Security.Secret+password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func (m *UsersManager) compareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(m.Config.Security.Secret+password))
}
