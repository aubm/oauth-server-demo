package security

import (
	"database/sql"
	"errors"
	"time"

	"github.com/RangelReale/osin"
)

type Storage struct {
	DB *sql.DB `inject:""`
}

func (s *Storage) Clone() osin.Storage {
	return s
}

func (s *Storage) Close() {
}

func (s *Storage) GetClient(id string) (osin.Client, error) {
	rows, err := s.DB.Query(`SELECT id, secret, redirect_uri
	FROM clients
	WHERE id = ?
	LIMIT 1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Client not found")
	}
	client := &Client{}
	if err := rows.Scan(&client.Id, &client.Secret, &client.RedirectUri); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) error {
	// Not supported
	return nil
}

func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	// Not supported
	return nil, nil
}

func (s *Storage) RemoveAuthorize(code string) error {
	// Not supported
	return nil
}

func (s *Storage) SaveAccess(data *osin.AccessData) error {
	stmt, err := s.DB.Prepare(`INSERT INTO access
	(access_token, refresh_token, expires_in, scope, redirect_uri, created_at, client_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		data.AccessToken,
		data.RefreshToken,
		data.ExpiresIn,
		data.Scope,
		data.RedirectUri,
		data.CreatedAt.Format("2006-01-02 15:04:05"),
		data.Client.GetId(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	rows, err := s.DB.Query(`SELECT c.id, c.secret, c.redirect_uri as client_redirect_uri,
	a.access_token, a.refresh_token, a.expires_in, a.scope, a.redirect_uri, a.created_at
	FROM access a
	INNER JOIN clients c ON (c.id = a.client_id)
	WHERE a.access_token = ?`, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Access not found")
	}
	d := &osin.AccessData{}
	c := &Client{}
	var createdAt string
	if err := rows.Scan(
		&c.Id, &c.Secret, &c.RedirectUri,
		&d.AccessToken, &d.RefreshToken, &d.ExpiresIn, &d.Scope, &d.RedirectUri, &createdAt,
	); err != nil {
		return nil, err
	}
	d.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	d.Client = c
	return d, nil
}

func (s *Storage) RemoveAccess(code string) error {
	stmt, err := s.DB.Prepare("DELETE FROM access WHERE access_token = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(code); err != nil {
		return err
	}
	return nil
}

func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	rows, err := s.DB.Query(`SELECT c.id, c.secret, c.redirect_uri as client_redirect_uri,
	a.access_token, a.refresh_token, a.expires_in, a.scope, a.redirect_uri, a.created_at
	FROM access a
	INNER JOIN clients c ON (c.id = a.client_id)
	WHERE a.refresh_token = ?`, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Access not found")
	}
	d := &osin.AccessData{}
	c := &Client{}
	var createdAt string
	if err := rows.Scan(
		&c.Id, &c.Secret, &c.RedirectUri,
		&d.AccessToken, &d.RefreshToken, &d.ExpiresIn, &d.Scope, &d.RedirectUri, &createdAt,
	); err != nil {
		return nil, err
	}
	d.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	d.Client = c
	return d, nil
}

func (s *Storage) RemoveRefresh(code string) error {
	stmt, err := s.DB.Prepare("DELETE FROM access WHERE refresh_token = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(code); err != nil {
		return err
	}
	return nil
}
