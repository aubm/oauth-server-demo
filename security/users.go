package security

import (
	"database/sql"

	"github.com/aubm/oauth-server-demo/config"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type NoUserFoundErr struct{}

func (_ NoUserFoundErr) Error() string {
	return "No user found"
}

type User struct {
	Id       string
	Email    string
	Password string
}

type UsersManager struct {
	DB          *sql.DB     `inject:""`
	Config      *config.App `inject:""`
	LoggerError interface {
		Printf(format string, v ...interface{})
	} `inject:"logger_error"`
}

func (m *UsersManager) Save(u User) error {
	if u.Id == "" {
		var err error
		u.Password, err = m.encrypt(u.Password)
		if err != nil {
			return err
		}
		u.Id = uuid.NewV4().String()
	}
	stmt, _ := m.DB.Prepare(`INSERT INTO users
	(id, email, password) VALUES (?, ?, ?)`)
	_, err := stmt.Exec(u.Id, u.Email, u.Password)
	if err != nil {
		m.LoggerError.Printf("failed to insert user into DB: %v", err)
	}
	return err
}

func (m *UsersManager) FindByCredentials(email, clearPassword string) (*User, error) {
	u, err := m.FindByEmail(email)
	if err != nil {
		return u, err
	}
	if err := m.compareHashAndPassword(u.Password, clearPassword); err != nil {
		return nil, NoUserFoundErr{}
	}
	return u, nil
}

func (m *UsersManager) FindByEmail(email string) (*User, error) {
	rows, err := m.DB.Query(`SELECT id, email, password
	FROM users
	WHERE email = ?
	LIMIT 1`, email)
	if err != nil {
		m.LoggerError.Printf("failed to load user from DB: %v", err)
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, NoUserFoundErr{}
	}

	u := new(User)
	if err := rows.Scan(&u.Id, &u.Email, &u.Password); err != nil {
		m.LoggerError.Printf("failed to scan user: %v", err)
		return nil, err
	}

	return u, nil
}

func (m *UsersManager) encrypt(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(m.Config.Security.Secret+password), bcrypt.DefaultCost)
	if err != nil {
		m.LoggerError.Printf("failed to encrypt password %v: %v", password, err)
		return "", err
	}
	return string(b[:]), nil
}

func (m *UsersManager) compareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(m.Config.Security.Secret+password))
}
