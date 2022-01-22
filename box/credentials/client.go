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

type client struct {
	clientID     string
	secret       string
	user         string
	userID       string
	enterpriseID string
}

func (c *client) Authenticate() (*AccessToken, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	form := url.Values{
		"client_id":        []string{c.clientID},
		"client_secret":    []string{c.secret},
		"grant_type":       []string{"client_credentials"},
		"box_subject_type": []string{c.user},
		"box_subject_id":   []string{c.enterpriseID},
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

func (c *client) load(bytes []byte) error {
	credentials := struct {
		Client struct {
			ClientID     string `json:"client-id"`
			Secret       string `json:"secret"`
			User         string `json:"user"`
			UserID       string `json:"user-id"`
			EnterpriseID string `json:"enterprise-id"`
		} `json:"client"`
	}{}

	if err := json.Unmarshal(bytes, &credentials); err != nil {
		return err
	}

	if credentials.Client.ClientID == "" {
		return fmt.Errorf("Invalid client ID (%v)", credentials.Client.ClientID)
	}

	if credentials.Client.Secret == "" {
		return fmt.Errorf("Invalid secret (%v)", credentials.Client.Secret)
	}

	if credentials.Client.User == "" {
		return fmt.Errorf("Invalid user (%v)", credentials.Client.User)
	}

	if credentials.Client.UserID == "" {
		return fmt.Errorf("Invalid user ID (%v)", credentials.Client.UserID)
	}

	if credentials.Client.EnterpriseID == "" {
		return fmt.Errorf("Invalid enterprise ID (%v)", credentials.Client.EnterpriseID)
	}

	c.clientID = credentials.Client.ClientID
	c.secret = credentials.Client.Secret
	c.user = credentials.Client.User
	c.userID = credentials.Client.UserID
	c.enterpriseID = credentials.Client.EnterpriseID

	return nil
}
