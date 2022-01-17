package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TemplateKey string

type Schema struct {
	Key    TemplateKey `json:"templateKey"`
	Name   string      `json:"displayName"`
	Fields []Field     `json:"fields"`
}

type Field struct {
	Type        string   `json:"type"`
	Key         string   `json:"key"`
	Name        string   `json:"displayName"`
	Description string   `json:"description"`
	Options     []Option `json:"options,omitempty"`
}

type Option struct {
	Key string `json:"key"`
}

func listTemplates(token string) (map[string]TemplateKey, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/metadata_templates/enterprise")

	rq, err := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", auth)
	rq.Header.Set("Content-Type", "application/json")
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
		return nil, fmt.Errorf("error retrieving list of templates (%v)", response.Status)
	}

	reply := struct {
		Limit   int `json:"limit"`
		Entries []struct {
			ID   string `json:"id"`
			Name string `json:"displayName"`
			Key  string `json:"templateKey"`
		} `json:"entries"`
	}{}

	if err := json.Unmarshal(body, &reply); err != nil {
		return nil, err
	}

	templates := map[string]TemplateKey{}

	for _, e := range reply.Entries {
		templates[e.Name] = TemplateKey(e.Key)
	}

	return templates, nil
}

func createTemplate(name string, fields []Field, token string) (TemplateKey, error) {
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

func deleteTemplate(key TemplateKey, token string) error {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/metadata_templates/%v/%v/schema", "enterprise", key)

	rq, err := http.NewRequest("DELETE", uri, nil)
	rq.Header.Set("Authorization", auth)
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accepts", "application/json")

	response, err := client.Do(rq)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if _, err := io.ReadAll(response.Body); err != nil {
		return err
	}

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting template (%v)", response.Status)
	}

	return nil
}

func getTemplate(template TemplateKey, token string) (*Schema, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/metadata_templates/enterprise/%v/schema", template)

	rq, err := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", auth)
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accepts", "application/json")

	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", response.Status)
	}

	schema := Schema{}
	if err := json.Unmarshal(body, &schema); err != nil {
		return nil, err
	}

	return &schema, nil
}
