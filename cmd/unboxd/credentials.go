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
		Box struct {
			Client *credentials.Client `json:"client,omitempty"`
			JWT    *credentials.JWT    `json:"jwt,omitempty"`
		} `json:"box"`
	}{}

	if err := json.Unmarshal(bytes, &credentials); err != nil {
		return nil, err
	}

	if credentials.Box.Client != nil {
		return credentials.Box.Client, nil
	} else if credentials.Box.JWT != nil {
		return credentials.Box.JWT, nil
	}

	return nil, fmt.Errorf("no valid credentials in file %v", file)
}
