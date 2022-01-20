package credentials

import _ "embed"

import (
	"testing"
)

//go:embed config_test.json
var config []byte

func TestJWTDecrypt(t *testing.T) {
	t.Skip()

	jwt := jwt{}

	if err := jwt.load(config); err != nil {
		t.Fatalf("Error initialising JWT (%v)", err)
	}

	if err := jwt.decrypt(); err != nil {
		t.Fatalf("Error decrypting JWT private key (%v)", err)
	}
}
