package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type File struct {
	ID   uint64
	Name string
	Tags []string
}

const fetchSize = 500

func get(fileID uint64, token string) (*File, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/files/%[1]v?fields=id,type,name,sha1,tags", fileID)

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

	if id, err := strconv.ParseUint(reply.ID, 10, 64); err != nil {
		return nil, err
	} else {
		return &File{
			ID:   id,
			Name: reply.Name,
			Tags: reply.Tags,
		}, nil
	}
}

func put(fileID uint64, content interface{}, token string) error {
	encoded, err := json.Marshal(content)
	if err != nil {
		return err
	}

	auth := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("https://api.box.com/2.0/files/%[1]v?fields=id,type,name,sha1,tags", fileID)

	rq, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(encoded))
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
