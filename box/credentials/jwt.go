package credentials

import (
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	JWT "github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/youmark/pkcs8"
)

type jwt struct {
	clientID     string
	secret       string
	publicKeyID  string
	privateKey   string
	passphrase   string
	enterpriseID string
}

type claims struct {
	JWT.RegisteredClaims
	BoxSubType string `json:"box_sub_type,omitempty"`
}

func (j *jwt) Authenticate() (*AccessToken, error) {
	pk, err := j.decrypt()
	if err != nil {
		return nil, err
	}

	token, err := j.assert(pk)
	if err != nil {
		return nil, err
	}

	return j.authenticate(token)
}

func (j *jwt) load(bytes []byte) error {
	credentials := struct {
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
	}{}

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

	j.clientID = credentials.BoxAppSettings.ClientID
	j.secret = credentials.BoxAppSettings.Secret
	j.publicKeyID = credentials.BoxAppSettings.AppAuth.PublicKeyID
	j.privateKey = credentials.BoxAppSettings.AppAuth.PrivateKey
	j.passphrase = credentials.BoxAppSettings.AppAuth.Passphrase
	j.enterpriseID = credentials.EnterpriseID

	return nil
}

func (j *jwt) decrypt() (*rsa.PrivateKey, error) {
	pwd := []byte(j.passphrase)
	block, _ := pem.Decode([]byte(j.privateKey))
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

func (j *jwt) assert(pk *rsa.PrivateKey) (*JWT.Token, error) {
	UUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	claims := claims{
		RegisteredClaims: JWT.RegisteredClaims{
			ID:        fmt.Sprintf("%v", UUID),
			Audience:  []string{"https://api.box.com/oauth2/token"},
			Issuer:    j.clientID,
			Subject:   j.enterpriseID,
			ExpiresAt: JWT.NewNumericDate(time.Now().Add(60 * time.Second)),
		},
		BoxSubType: "enterprise",
	}

	signer, err := JWT.NewSignerRS(JWT.RS512, pk)
	if err != nil {
		return nil, err
	}

	return JWT.NewBuilder(signer, JWT.WithKeyID(j.publicKeyID)).Build(claims)
}

func (j *jwt) authenticate(t *JWT.Token) (*AccessToken, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	assertion := string(t.Bytes())

	form := url.Values{
		"client_id":     []string{j.clientID},
		"client_secret": []string{j.secret},
		"grant_type":    []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":     []string{assertion},
	}

	rq, err := http.NewRequest("POST", "https://api.box.com/oauth2/token", strings.NewReader(form.Encode()))
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rq.Header.Set("Accepts", "application/json")

	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authorization request failed (%s)", response.Status)
	}

	token := struct {
		AccessToken  string   `json:"access_token"`
		ExpiresIn    int      `json:"expires_in"`
		RestrictedTo []string `json:"restricted_to"`
		TokenType    string   `json:"token_type"`
	}{}

	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &AccessToken{
		Token:  token.AccessToken,
		Expiry: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}, nil
}
