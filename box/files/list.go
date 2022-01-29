package files

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func List(folderID string, token string) ([]File, error) {
	files := []File{}
	auth := fmt.Sprintf("Bearer %s", token)
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	uri := fmt.Sprintf("https://api.box.com/2.0/folders/%[1]v/items?fields=id,type,name,sha1,tags&limit=%[2]v&usemarker=true", folderID, fetchSize)

	for {
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
			return nil, fmt.Errorf("error retrieving list of files (%v)", response.Status)
		}

		reply := struct {
			TotalCount int `json:"total_count"`
			Entries    []struct {
				Type string   `json:"type"`
				ID   string   `json:"id"`
				Name string   `json:"name"`
				Tags []string `json:"tags"`
			} `json:"entries"`
			NextMarker string `json:"next_marker,omitempty"`
		}{}

		if err := json.Unmarshal(body, &reply); err != nil {
			return nil, err
		}

		for _, e := range reply.Entries {
			if e.Type == "file" {
				files = append(files, File{
					ID:       e.ID,
					Filename: e.Name,
					Tags:     e.Tags,
				})
			}
		}

		if reply.NextMarker == "" {
			break
		}

		uri = fmt.Sprintf("https://api.box.com/2.0/folders/%[1]v/items?fields=id,type,name,sha1&limit=%[2]v&marker=%[3]v&usemarker=true", folderID, fetchSize,reply.NextMarker)
	}

	return files, nil
}
