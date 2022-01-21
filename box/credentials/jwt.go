package credentials

import (
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/youmark/pkcs8"
)

type jwtx struct {
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

type claims struct {
	jwt.RegisteredClaims
	BoxSubType string `json:"box_sub_type,omitempty"`
}

func (j *jwtx) Authenticate() (*AccessToken, error) {
	_, err := j.decrypt()
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (j *jwtx) load(bytes []byte) error {
	credentials := jwtx{}
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

func (j *jwtx) authenticate() (*AccessToken, error) {
	pk, err := j.decrypt()
	if err != nil {
		return nil, err
	}

	UUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%v", UUID),
			Audience:  []string{"https://api.box.com/oauth2/token"},
			Issuer:    j.BoxAppSettings.ClientID,
			Subject:   j.EnterpriseID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Second)),
		},
		BoxSubType: "enterprise",
	}

	signer, err := jwt.NewSignerRS(jwt.RS512, pk)
	if err != nil {
		return nil, err
	}

	token, err := jwt.NewBuilder(signer, jwt.WithKeyID(j.BoxAppSettings.AppAuth.PublicKeyID)).Build(claims)
	if err != nil {
		return nil, err
	}

	fmt.Printf(">>>>>>>>> algorithm: %#v\n", token.Header())
	fmt.Printf(">>>>>>>>> claims:    %v\n", string(token.Claims()))

	return nil, nil
}

func (j *jwtx) decrypt() (*rsa.PrivateKey, error) {
	pwd := []byte(j.BoxAppSettings.AppAuth.Passphrase)
	block, _ := pem.Decode([]byte(j.BoxAppSettings.AppAuth.PrivateKey))
	if block == nil || block.Type != "ENCRYPTED PRIVATE KEY" {
		return nil, fmt.Errorf("Invalid private key")
	}

	key, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, pwd)
	if err != nil {
		return nil, err
	}

	if pk, ok := key.(*rsa.PrivateKey); !ok {
		return nil, fmt.Errorf("Invalid private key")
	} else {
		return pk, nil
	}
}
