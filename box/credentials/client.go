package credentials

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	clientID     string
	secret       string
	user         string
	userID       string
	enterpriseID string
}

func (c *Client) Authenticate() (*AccessToken, error) {
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

func (c Client) Hash() string {
	s := fmt.Sprintf("%v:%v:%v", c.clientID, c.userID, c.enterpriseID)
	hash := sha256.Sum256([]byte(s))

	return fmt.Sprintf("%x", hash)
}

func (c *Client) UnmarshalJSON(bytes []byte) error {
	credentials := struct {
		ClientID     string `json:"client-id"`
		Secret       string `json:"secret"`
		User         string `json:"user"`
		UserID       string `json:"user-id"`
		EnterpriseID string `json:"enterprise-id"`
	}{}

	if err := json.Unmarshal(bytes, &credentials); err != nil {
		return err
	}

	c.clientID = credentials.ClientID
	c.secret = credentials.Secret
	c.user = credentials.User
	c.userID = credentials.UserID
	c.enterpriseID = credentials.EnterpriseID

	return nil
}
