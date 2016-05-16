package security

import (
	"encoding/json"
	"time"

	"github.com/RangelReale/osin"
	"gopkg.in/redis.v3"
)

type AccessDataManager struct {
	SetGetDel interface {
		Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
		Get(key string) *redis.StringCmd
		Del(keys ...string) *redis.IntCmd
	} `inject:""`
}

func (m *AccessDataManager) Save(data *osin.AccessData, accessDuration, refreshDuration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if data.RefreshToken != "" {
		if err := m.SetGetDel.Set("refresh_"+data.RefreshToken, b, refreshDuration).Err(); err != nil {
			return err
		}
	}
	return m.SetGetDel.Set("access_"+data.AccessToken, b, accessDuration).Err()
}

func (m *AccessDataManager) DeleteAccess(code string) error {
	return m.SetGetDel.Del("access_" + code).Err()
}

func (m *AccessDataManager) DeleteRefresh(code string) error {
	return m.SetGetDel.Del("refresh_" + code).Err()
}

func (m *AccessDataManager) FindAccess(code string) (*osin.AccessData, error) {
	return m.find("access_" + code)
}

func (m *AccessDataManager) FindRefresh(code string) (*osin.AccessData, error) {
	return m.find("refresh_" + code)
}

func (m *AccessDataManager) find(code string) (*osin.AccessData, error) {
	b, err := m.SetGetDel.Get(code).Bytes()
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	access := &osin.AccessData{}
	if cData, ok := data["Client"].(map[string]interface{}); ok {
		c := &Client{}
		if v, ok := cData["Id"].(string); ok {
			c.Id = v
		}
		if v, ok := cData["Secret"].(string); ok {
			c.Secret = v
		}
		if v, ok := cData["RedirectUri"].(string); ok {
			c.RedirectUri = v
		}
		if v, ok := cData["UserData"]; ok {
			c.UserData = v
		}
		access.Client = c
	}
	if v, ok := data["AccessToken"].(string); ok {
		access.AccessToken = v
	}
	if v, ok := data["RefreshToken"].(string); ok {
		access.RefreshToken = v
	}
	if v, ok := data["ExpiresIn"].(int32); ok {
		access.ExpiresIn = v
	}
	if v, ok := data["Scope"].(string); ok {
		access.Scope = v
	}
	if v, ok := data["RedirectUri"].(string); ok {
		access.RedirectUri = v
	}
	if v, ok := data["CreatedAt"].(string); ok {
		date, _ := time.Parse(time.RFC3339, v)
		access.CreatedAt = date
	}
	if uData, ok := data["UserData"].(map[string]interface{}); ok {
		u := &User{}
		if v, ok := uData["Id"].(string); ok {
			u.Id = v
		}
		if v, ok := uData["Email"].(string); ok {
			u.Email = v
		}
		if v, ok := uData["Password"].(string); ok {
			u.Password = v
		}
		access.UserData = u
	}

	return access, nil
}
