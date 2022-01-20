package credentials

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/youmark/pkcs8"
)

type jwt struct {
	BoxAppSettings struct {
		ClientID string `json:"clientID"`
		Secret   string `json:"clientSecret"`
		AppAuth  struct {
			PublicKeyID string `json:"publicKeyID"`
			PrivateKey  string `json:"privateKey"`
			Passphrase  string `json:"passphrase"`
		} `json:"appAuth"`
	} `json:"boxAppSettings"`

	EnterpriseID string `json:"enterpriseID"`
}

func (j *jwt) Authenticate() (*AccessToken, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (j *jwt) load(bytes []byte) error {
	credentials := jwt{}
	if err := json.Unmarshal(bytes, &credentials); err != nil {
		return err
	}

	if credentials.BoxAppSettings.ClientID == "" {
		return fmt.Errorf("Invalid client ID (%v)", credentials.BoxAppSettings.ClientID)
	}

	if credentials.BoxAppSettings.Secret == "" {
		return fmt.Errorf("Invalid secret (%v)", credentials.BoxAppSettings.Secret)
	}

	if credentials.BoxAppSettings.AppAuth.PublicKeyID == "" {
		return fmt.Errorf("Invalid public key ID (%v)", credentials.BoxAppSettings.AppAuth.PublicKeyID)
	}

	if credentials.BoxAppSettings.AppAuth.PrivateKey == "" {
		return fmt.Errorf("Invalid private key (%v)", credentials.BoxAppSettings.AppAuth.PrivateKey)
	}

	if credentials.EnterpriseID == "" {
		return fmt.Errorf("Invalid enterprise ID (%v)", credentials.EnterpriseID)
	}

	j.BoxAppSettings.ClientID = credentials.BoxAppSettings.ClientID
	j.BoxAppSettings.Secret = credentials.BoxAppSettings.Secret
	j.BoxAppSettings.AppAuth.PublicKeyID = credentials.BoxAppSettings.AppAuth.PublicKeyID
	j.BoxAppSettings.AppAuth.PrivateKey = credentials.BoxAppSettings.AppAuth.PrivateKey
	j.BoxAppSettings.AppAuth.Passphrase = credentials.BoxAppSettings.AppAuth.Passphrase
	j.EnterpriseID = credentials.EnterpriseID

	return nil
}

func (j *jwt) decrypt() error {
	pwd := []byte(j.BoxAppSettings.AppAuth.Passphrase)
	block, _ := pem.Decode([]byte(j.BoxAppSettings.AppAuth.PrivateKey))
	if block == nil || block.Type != "ENCRYPTED PRIVATE KEY" {
		return fmt.Errorf("Invalid private key")
	}

	key, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, pwd)
	fmt.Printf(">>> KEY: %v\n", key)
	fmt.Printf(">>> ERROR: %v\n", err)

	return fmt.Errorf("NOT IMPLEMENTED")
}
