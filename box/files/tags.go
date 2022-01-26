package files

import (
	"fmt"
	"sort"
)

func Tag(fileID string, tag string, token string) error {
	file, err := get(fileID, token)
	if err != nil {
		return err
	} else if file == nil {
		return fmt.Errorf("invalid file returned for %v", fileID)
	}

	tags := []string{}
	for _, t := range file.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}

	tags = append(tags, tag)
	if equal(tags, file.Tags) {
		return nil
	}

	info := struct {
		Tags []string `json:"tags"`
	}{
		Tags: tags,
	}

	return put(fileID, info, token)
}

func Untag(fileID string, tag string, token string) error {
	file, err := get(fileID, token)
	if err != nil {
		return err
	} else if file == nil {
		return fmt.Errorf("invalid file returned for %v", fileID)
	}

	tags := []string{}
	for _, t := range file.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}

	if equal(tags, file.Tags) {
		return nil
	}

	info := struct {
		Tags []string `json:"tags"`
	}{
		Tags: tags,
	}

	return put(fileID, info, token)
}

func Retag(fileID string, oldTag, newTag string, token string) error {
	file, err := get(fileID, token)
	if err != nil {
		return err
	} else if file == nil {
		return fmt.Errorf("invalid file returned for %v", fileID)
	}

	tags := []string{}
	for _, t := range file.Tags {
		if t != oldTag {
			tags = append(tags, t)
		}
	}

	if equal(tags, file.Tags) {
		return fmt.Errorf("file %v does not have tag '%v'", fileID, oldTag)
	}

	for _, t := range file.Tags {
		if t == newTag {
			return fmt.Errorf("file %v already has tag '%v'", fileID, newTag)
		}
	}

	tags = append(tags, newTag)
	if equal(tags, file.Tags) {
		return nil
	}

	info := struct {
		Tags []string `json:"tags"`
	}{
		Tags: tags,
	}

	return put(fileID, info, token)
}

func equal(p, q []string) bool {
	if len(p) != len(q) {
		return false
	}

	sort.Strings(p)
	sort.Strings(q)

	for i, u := range p {
		if u != q[i] {
			return false
		}
	}

	return true
}
