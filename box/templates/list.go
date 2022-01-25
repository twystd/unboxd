package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func ListTemplates(token string) (map[string]TemplateKey, error) {
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
