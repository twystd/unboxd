package templates

import ()

type TemplateKey string

type Schema struct {
	Key    TemplateKey `json:"templateKey"`
	Name   string      `json:"displayName"`
	Fields []Field     `json:"fields"`
}

type Field struct {
	Type        string   `json:"type"`
	Key         string   `json:"key"`
	Name        string   `json:"displayName"`
	Description string   `json:"description"`
	Options     []Option `json:"options,omitempty"`
}

type Option struct {
	Key string `json:"key"`
}
