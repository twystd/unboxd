package commands

import (
	"encoding/json"
	"errors"
	"io/fs"
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

func resume(file string) ([]uint64, []folder, error) {
	checkpoint := struct {
		Queue   []uint64 `json:"queue"`
		Folders []folder `json:"folders"`
	}{}

	if _, err := os.Stat(file); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, nil, err
	} else if err != nil {
		return []uint64{}, []folder{}, nil
	}

	if bytes, err := os.ReadFile(file); err != nil {
		return nil, nil, err
	} else if err := json.Unmarshal(bytes, &checkpoint); err != nil {
		return nil, nil, err
	} else {
		return checkpoint.Queue, checkpoint.Folders, nil
	}
}
