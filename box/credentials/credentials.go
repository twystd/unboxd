package credentials

import (
	"fmt"
	"os"
	"time"
)

type Credentials interface {
	Authenticate() (*AccessToken, error)
}

type AccessToken struct {
	Token  string
	Expiry time.Time
}

func (t AccessToken) IsValid() bool {
	renew := time.Now().Add(10 * time.Minute)

	return t.Token != "" && t.Expiry.After(renew)
}

func NewCredentials(file string) (Credentials, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	a := app{}
	if err := a.load(bytes); err == nil {
		return &a, nil
	}

	j := jwt{}
	if err := j.load(bytes); err == nil {
		return &j, nil
	}

	return nil, fmt.Errorf("%v - invalid credentials", file)
}
