package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func checkpoint(file string, queue []uint64, folders []folder) error {
	checkpoint := struct {
		Queue   []uint64 `json:"queue"`
		Folders []folder `json:"folders"`
	}{
		Queue:   queue,
		Folders: folders,
	}

	if err := os.MkdirAll(filepath.Dir(file), 0750); err != nil && !os.IsExist(err) {
		return err
	}

	if bytes, err := json.MarshalIndent(checkpoint, "", "  "); err != nil {
		return err
	} else if err := os.WriteFile(file, bytes, 0666); err != nil {
		return err
	}

	return nil
}
