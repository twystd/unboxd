package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FileID string

type File struct {
	ID       FileID
	Filename string
	Tags     []string
}

func TagFile(fileID FileID, tag string, token string) error {
	file, err := get(fileID, token)
	if err != nil {
		return err
	} else if file == nil {
		return fmt.Errorf("invalid file returned for %v", fileID)
	}

	for _, t := range file.Tags {
		if t == tag {
			return nil
		}
	}

	info := struct {
		Tags []string `json:"tags"`
	}{
		Tags: append(file.Tags, tag),
	}

	return put(fileID, info, token)
}

func get(fileID FileID, token string) (*File, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/files/%[1]v?fields=id,type,name,sha1,tags", fileID)

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
		return nil, fmt.Errorf("%v: error retrieving file information (%v)", fileID, response.Status)
	}

	reply := struct {
		Type string   `json:"type"`
		ID   string   `json:"id"`
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{}

	if err := json.Unmarshal(body, &reply); err != nil {
		return nil, err
	}

	return &File{
		ID:       FileID(reply.ID),
		Filename: reply.Name,
		Tags:     reply.Tags,
	}, nil
}

func put(fileID FileID, content interface{}, token string) error {
	encoded, err := json.Marshal(content)
	if err != nil {
		return err
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/files/%[1]v?fields=id,type,name,sha1,tags", fileID)

	rq, err := http.NewRequest("PUT", uri, bytes.NewBuffer(encoded))
	rq.Header.Set("Authorization", auth)
	rq.Header.Set("Content-Type", "application/json")
	rq.Header.Set("Accepts", "application/json")

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	response, err := client.Do(rq)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if _, err := io.ReadAll(response.Body); err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error tagging file (%v)", response.Status)
	}

	return nil
}
