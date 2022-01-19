package credentials

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type app struct {
	ClientID     string `json:"client-id"`
	Secret       string `json:"secret"`
	User         string `json:"user"`
	UserID       string `json:"user-id"`
	EnterpriseID string `json:"enterprise-id"`
}

func (a *app) Authenticate() (*AccessToken, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	form := url.Values{
		"client_id":        []string{a.ClientID},
		"client_secret":    []string{a.Secret},
		"grant_type":       []string{"client_credentials"},
		"box_subject_type": []string{a.User},
		"box_subject_id":   []string{a.EnterpriseID},
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

func (a *app) load(bytes []byte) error {
	credentials := app{}
	if err := json.Unmarshal(bytes, &credentials); err != nil {
		return err
	}

	if credentials.ClientID == "" {
		return fmt.Errorf("Invalid client ID (%v)", credentials.ClientID)
	}

	if credentials.Secret == "" {
		return fmt.Errorf("Invalid secret (%v)", credentials.Secret)
	}

	if credentials.User == "" {
		return fmt.Errorf("Invalid user (%v)", credentials.User)
	}

	if credentials.UserID == "" {
		return fmt.Errorf("Invalid user ID (%v)", credentials.UserID)
	}

	if credentials.EnterpriseID == "" {
		return fmt.Errorf("Invalid enterprise ID (%v)", credentials.EnterpriseID)
	}

	a.ClientID = credentials.ClientID
	a.Secret = credentials.Secret
	a.User = credentials.User
	a.UserID = credentials.UserID
	a.EnterpriseID = credentials.EnterpriseID

	return nil
}
