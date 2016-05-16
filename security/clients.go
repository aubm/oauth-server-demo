package security

import (
	"database/sql"
	"errors"
)

type Client struct {
	Id          string
	Secret      string
	RedirectUri string
	UserData    interface{}
}

func (d *Client) GetId() string {
	return d.Id
}

func (d *Client) GetSecret() string {
	return d.Secret
}

func (d *Client) GetRedirectUri() string {
	return d.RedirectUri
}

func (d *Client) GetUserData() interface{} {
	return d.UserData
}

type ClientsManager struct {
	DB *sql.DB `inject:""`
}

func (m *ClientsManager) FindOne(id string) (*Client, error) {
	rows, err := m.DB.Query(`SELECT id, secret, redirect_uri
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
