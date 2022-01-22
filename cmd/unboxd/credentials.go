package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/twystd/unboxd/box/credentials"
)

func NewCredentials(file string) (credentials.Credentials, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	credentials := struct {
		Client *credentials.Client `json:"client,omitempty"`
		JWT    *credentials.JWT    `json:"jwt,omitempty"`
	}{}

	if err := json.Unmarshal(bytes, &credentials); err != nil {
		return nil, err
	}

	if credentials.Client != nil {
		return credentials.Client, nil
	} else if credentials.JWT != nil {
		return credentials.JWT, nil
	}

	return nil, fmt.Errorf("No valid credentials in file %v", file)
}
