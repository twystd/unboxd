package box

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessToken struct {
	token  string
	expiry time.Time
}

func (t AccessToken) IsValid() bool {
	renew := time.Now().Add(10 * time.Minute)

	return t.token != "" && t.expiry.After(renew)
}

func authenticate(ID, secret, user, userID string) (*AccessToken, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	form := url.Values{
		"client_id":        []string{ID},
		"client_secret":    []string{secret},
		"grant_type":       []string{"client_credentials"},
		"box_subject_type": []string{user},
		"box_subject_id":   []string{userID},
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
		token:  token.AccessToken,
		expiry: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}, nil
}
