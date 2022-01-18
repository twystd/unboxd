package credentials

import (
	"encoding/json"
	"os"
)

type Credentials struct {
	ClientID string  `json:"client-id"`
	Secret   string  `json:"secret"`
	User     string  `json:"user"`
	UserID   string  `json:"user-id"`
	Folders  Folders `json:"folders"`
}

type Folders struct {
	Photos  string `json:"photos"`
	Pending string `json:"pending"`
}

func (c *Credentials) Load(file string) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, c)
}
