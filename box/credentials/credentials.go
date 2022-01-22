package credentials

import (
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
