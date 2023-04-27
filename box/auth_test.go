package box

import _ "embed"

import (
	"encoding/json"
	"testing"
)

//go:embed auth_test.json
var config []byte

func TestJWTAuthenticate(t *testing.T) {
	t.Skip()

	j := JWT{}

	if err := json.Unmarshal(config, &j); err != nil {
		t.Fatalf("Error initialising JWT (%v)", err)
	}

	token, err := j.Authenticate()
	if err != nil {
		t.Fatalf("Error decrypting JWT private key (%v)", err)
	} else if token == nil {
		t.Errorf("Invalid Box access token (%v)", token)
	}
}
