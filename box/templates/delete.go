package templates

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Delete(key TemplateKey, token string) error {
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
