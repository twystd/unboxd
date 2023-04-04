package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Get(template TemplateKey, token string) (*Schema, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/metadata_templates/enterprise/%v/schema", template)

	rq, _ := http.NewRequest("GET", uri, nil)
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
