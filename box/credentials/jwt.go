package credentials

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/youmark/pkcs8"
)

type JWT struct {
	clientID     string
	secret       string
	publicKeyID  string
	privateKey   string
	passphrase   string
	enterpriseID string
}

type claims struct {
	jwt.RegisteredClaims
	BoxSubType string `json:"box_sub_type,omitempty"`
}

func (j *JWT) Authenticate() (*AccessToken, error) {
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

func (j JWT) Hash() string {
	s := fmt.Sprintf("%v:%v:%v", j.clientID, j.publicKeyID, j.enterpriseID)
	hash := sha256.Sum256([]byte(s))

	return fmt.Sprintf("%x", hash)
}

func (j *JWT) UnmarshalJSON(bytes []byte) error {
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

	j.clientID = credentials.BoxAppSettings.ClientID
	j.secret = credentials.BoxAppSettings.Secret
	j.publicKeyID = credentials.BoxAppSettings.AppAuth.PublicKeyID
	j.privateKey = credentials.BoxAppSettings.AppAuth.PrivateKey
	j.passphrase = credentials.BoxAppSettings.AppAuth.Passphrase
	j.enterpriseID = credentials.EnterpriseID

	return nil
}

func (j *JWT) decrypt() (*rsa.PrivateKey, error) {
	pwd := []byte(j.passphrase)
	block, _ := pem.Decode([]byte(j.privateKey))
	if block == nil || block.Type != "ENCRYPTED PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key")
	}

	key, err := pkcs8.ParsePKCS8PrivateKey(block.Bytes, pwd)
	if err != nil {
		return nil, err
	}

	if pk, ok := key.(*rsa.PrivateKey); !ok {
		return nil, fmt.Errorf("invalid private key")
	} else {
		return pk, nil
	}
}

func (j *JWT) assert(pk *rsa.PrivateKey) (*jwt.Token, error) {
	UUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%v", UUID),
			Audience:  []string{"https://api.box.com/oauth2/token"},
			Issuer:    j.clientID,
			Subject:   j.enterpriseID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Second)),
		},
		BoxSubType: "enterprise",
	}

	signer, err := jwt.NewSignerRS(jwt.RS512, pk)
	if err != nil {
		return nil, err
	}

	return jwt.NewBuilder(signer, jwt.WithKeyID(j.publicKeyID)).Build(claims)
}

func (j *JWT) authenticate(t *jwt.Token) (*AccessToken, error) {
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

	rq, _ := http.NewRequest("POST", "https://api.box.com/oauth2/token", strings.NewReader(form.Encode()))
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
