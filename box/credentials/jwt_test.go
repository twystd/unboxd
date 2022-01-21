package credentials

import _ "embed"

import (
	"testing"
)

//go:embed config_test.json
var config []byte

func TestJWTDecrypt(t *testing.T) {
	t.Skip()

	j := jwtx{}

	if err := j.load(config); err != nil {
		t.Fatalf("Error initialising JWT (%v)", err)
	}

	if key, err := j.decrypt(); err != nil {
		t.Fatalf("Error decrypting JWT private key (%v)", err)
	} else if key == nil {
		t.Errorf("Invalid private key (%v)", key)
	}
}

func TestJWTAuthenticate(t *testing.T) {
	t.Skip()

	j := jwtx{}

	if err := j.load(config); err != nil {
		t.Fatalf("Error initialising JWT (%v)", err)
	}

	if token, err := j.authenticate(); err != nil {
		t.Fatalf("Error authenticating to Box (%v)", err)
	} else if token == nil {
		t.Errorf("Invalid auth token (%v)", token)
	}
}
