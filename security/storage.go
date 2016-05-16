package security

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RangelReale/osin"
)

type Storage struct {
	DB                *sql.DB            `inject:""`
	AccessDataManager *AccessDataManager `inject:""`
	ClientsManager    interface {
		FindOne(id string) (*Client, error)
	} `inject:""`
}

func (s *Storage) Clone() osin.Storage {
	return s
}

func (s *Storage) Close() {
}

func (s *Storage) GetClient(id string) (osin.Client, error) {
	return s.ClientsManager.FindOne(id)
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
	accessDuration, _ := time.ParseDuration(fmt.Sprintf("%vs", data.ExpiresIn))
	refreshDuration, _ := time.ParseDuration("168h") // 1 week
	return s.AccessDataManager.Save(data, accessDuration, refreshDuration)
}

func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	return s.AccessDataManager.FindAccess(code)
}

func (s *Storage) RemoveAccess(code string) error {
	return s.AccessDataManager.DeleteAccess(code)
}

func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	return s.AccessDataManager.FindRefresh(code)
}

func (s *Storage) RemoveRefresh(code string) error {
	return s.AccessDataManager.DeleteRefresh(code)
}
