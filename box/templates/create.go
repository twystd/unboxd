package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func CreateTemplate(name string, fields []Field, token string) (TemplateKey, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/metadata_templates/schema")

	template := struct {
		Scope  string  `json:"scope"`
		Name   string  `json:"displayName"`
		Fields []Field `json:"fields"`
	}{
		Scope:  "enterprise",
		Name:   name,
		Fields: fields,
	}

	encoded, err := json.Marshal(template)
	if err != nil {
		return "", err
	}

	rq, err := http.NewRequest("POST", uri, bytes.NewBuffer(encoded))
	rq.Header.Set("Authorization", auth)
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accepts", "application/json")

	response, err := client.Do(rq)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("error creating template (%v)", response.Status)
	}

	reply := struct {
		ID  string `json:"id"`
		Key string `json:"templateKey"`
	}{}

	if err := json.Unmarshal(body, &reply); err != nil {
		return "", err
	}

	return TemplateKey(reply.Key), nil
}
