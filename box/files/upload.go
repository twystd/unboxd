package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Upload(file string, folder string, token string) (string, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	filename := filepath.Base(file)
	attributes := struct {
		Name   string `json:"name"`
		Parent struct {
			ID string `json:"id"`
		} `json:"parent"`
	}{
		Name: filename,
		Parent: struct {
			ID string `json:"id"`
		}{
			ID: folder,
		},
	}

	a, err := json.Marshal(attributes)
	if err != nil {
		return "", err
	}

	r, err := os.Open(file)
	if err != nil {
		return "", err
	} else {
		defer r.Close()
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err = writer.WriteField("attributes", string(a)); err != nil {
		return "", err
	}

	if part, err := writer.CreateFormFile("file", filename); err != nil {
		return "", err
	} else if _, err = io.Copy(part, r); err != nil {
		return "", err
	}

	if err = writer.Close(); err != nil {
		return "", err
	}

	rq, err := http.NewRequest("POST", "https://upload.box.com/api/2.0/files/content", body)
	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rq.Header.Set("Accepts", "application/json")
	rq.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := client.Do(rq)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	reply, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// ... ok - acknowledge with Box file ID
	if response.StatusCode == http.StatusCreated {
		info := struct {
			TotalCount int `json:"total_count"`
			Entries    []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Size uint64 `json:"size"`
				SHA1 string `json:"sha1"`
			} `json:"entries"`
		}{}

		if err := json.Unmarshal(reply, &info); err != nil {
			return "", err
		} else if info.TotalCount != 1 {
			return "", fmt.Errorf("invalid response - total count:%v", info.TotalCount)
		} else if len(info.Entries) != 1 {
			return "", fmt.Errorf("invalid response - entries:%v", len(info.Entries))
		}

		return info.Entries[0].ID, nil
	}

	return "", fmt.Errorf("upload request failed (%s)", response.Status)
}
