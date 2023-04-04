package folders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func List(folderID uint64, token string) ([]Folder, error) {
	folders := []Folder{}
	auth := fmt.Sprintf("Bearer %s", token)
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	uri := fmt.Sprintf("https://api.box.com/2.0/folders/%[1]v/items?fields=id,type,name,sha1&limit=%[2]v&usemarker=true", folderID, fetchSize)

	for {
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
		}

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error retrieving list of folders (%v)", response.Status)
		}

		reply := struct {
			TotalCount int `json:"total_count"`
			Entries    []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"entries"`
			NextMarker string `json:"next_marker,omitempty"`
		}{}

		if err := json.Unmarshal(body, &reply); err != nil {
			return nil, err
		}

		for _, e := range reply.Entries {
			if e.Type == "folder" {
				if id, err := strconv.ParseUint(e.ID, 10, 64); err == nil {
					folders = append(folders, Folder{
						ID:   id,
						Name: e.Name,
					})
				}
			}
		}

		fmt.Printf(">> folder:%v  total:%-4v entries:%-4v  folders:%-4v\n", folderID, reply.TotalCount, len(reply.Entries), len(folders))

		if reply.NextMarker == "" {
			break
		}

		uri = fmt.Sprintf("https://api.box.com/2.0/folders/%[1]v/items?fields=id,type,name,sha1&limit=%[2]v&marker=%[3]v&usemarker=true", folderID, fetchSize, reply.NextMarker)
	}

	return folders, nil
}
