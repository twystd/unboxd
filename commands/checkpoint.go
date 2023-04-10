package commands

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type Checkpoint struct {
	Queue   []QueueItem `json:"queue"`
	Folders []folder    `json:"folders"`
}

type QueueItem struct {
	ID   uint64 `json:"ID"`
	Path string `json:"path"`
}

func checkpoint(file string, queue []QueueItem, folders []folder) error {
	checkpoint := Checkpoint{
		Queue:   queue,
		Folders: folders,
	}

	if file != "" {
		if err := os.MkdirAll(filepath.Dir(file), 0750); err != nil && !os.IsExist(err) {
			return err
		}

		if bytes, err := json.MarshalIndent(checkpoint, "", "  "); err != nil {
			return err
		} else if err := os.WriteFile(file, bytes, 0666); err != nil {
			return err
		}
	}

	return nil
}

func resume(file string, restart bool) ([]QueueItem, []folder, error) {
	checkpoint := Checkpoint{}

	if file != "" && !restart {
		if _, err := os.Stat(file); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, nil, err
		} else if err != nil {
			return []QueueItem{}, []folder{}, nil
		}

		if bytes, err := os.ReadFile(file); err != nil {
			return nil, nil, err
		} else if err := json.Unmarshal(bytes, &checkpoint); err != nil {
			return nil, nil, err
		} else {
			return checkpoint.Queue, checkpoint.Folders, nil
		}
	}

	return []QueueItem{}, []folder{}, nil
}
