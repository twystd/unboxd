package commands

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type Checkpoint struct {
	Hash    string      `json:"id"`
	Queue   []QueueItem `json:"queue"`
	Folders []folder    `json:"folders"`
	Files   []file      `json:"files"`
}

type QueueItem struct {
	ID   uint64 `json:"ID"`
	Path string `json:"path"`
}

func checkpoint(file string, queue []QueueItem, folders []folder, files []file, hash string) error {
	checkpoint := Checkpoint{
		Hash:    hash,
		Queue:   queue,
		Folders: folders,
		Files:   files,
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

func resume(chkpt string, hash string, restart bool) ([]QueueItem, []folder, []file, error) {
	checkpoint := Checkpoint{}

	if chkpt != "" && !restart {
		if _, err := os.Stat(chkpt); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, nil, nil, err
		} else if err != nil {
			return []QueueItem{}, []folder{}, []file{}, nil
		}

		if bytes, err := os.ReadFile(chkpt); err != nil {
			return nil, nil, nil, err
		} else if err := json.Unmarshal(bytes, &checkpoint); err != nil {
			return nil, nil, nil, err
		} else if checkpoint.Hash != hash {
			return []QueueItem{}, []folder{}, []file{}, nil
		} else {
			return checkpoint.Queue, checkpoint.Folders, checkpoint.Files, nil
		}
	}

	return []QueueItem{}, []folder{}, []file{}, nil
}
