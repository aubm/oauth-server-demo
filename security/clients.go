package security

import "database/sql"

type NoClientFoundErr struct{}

func (_ NoClientFoundErr) Error() string {
	return "No client found"
}

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
	DB          *sql.DB `inject:""`
	LoggerError interface {
		Printf(format string, v ...interface{})
	} `inject:"logger_error"`
}

func (m *ClientsManager) FindOne(id string) (*Client, error) {
	rows, err := m.DB.Query(`SELECT id, secret, redirect_uri
	FROM clients
	WHERE id = ?
	LIMIT 1`, id)
	if err != nil {
		m.LoggerError.Printf("failed to load client from DB: %v", err)
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NoClientFoundErr{}
	}
	client := &Client{}
	if err := rows.Scan(&client.Id, &client.Secret, &client.RedirectUri); err != nil {
		m.LoggerError.Printf("failed to scan client: %v", err)
		return nil, err
	}
	return client, nil
}
