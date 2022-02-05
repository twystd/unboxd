package lib

import (
	"testing"
)

func TestGlobMatchAll(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", true},
		{"/alpha", true},
		{"/alpha/pending", true},
		{"/alpha/photos", true},
		{"/alpha/photos/new", true},
		{"/alpha/photos/new/today", true},
		{"/beta", true},
		{"/beta/pending", true},
		{"/beta/photos", true},
	}

	g := NewGlob("")

	for _, v := range tests {
		if match := g.Match(v.path); match != v.expected {
			t.Errorf("Incorrect match for '%s' - expected:%v, got:%v", v.path, v.expected, match)
		}
	}
}

func TestGlobMatchRoot(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", false},
		{"/alpha", true},
		{"/alpha/pending", false},
		{"/alpha/photos", false},
		{"/alpha/photos/new", false},
		{"/alpha/photos/new/today", false},
		{"/beta", true},
		{"/beta/pending", false},
		{"/beta/photos", false},
	}

	g := NewGlob("/")

	for _, v := range tests {
		if match := g.Match(v.path); match != v.expected {
			t.Errorf("Incorrect match for '%s' - expected:%v, got:%v", v.path, v.expected, match)
		}
	}
}

func TestGlobMatchExactPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", false},
		{"/alpha", true},
		{"/alpha/pending", false},
		{"/alpha/photos", false},
		{"/alpha/photos/new", false},
		{"/alpha/photos/new/today", false},
		{"/beta", false},
		{"/beta/pending", false},
		{"/beta/photos", false},
	}

	g := NewGlob("/alpha")

	for _, v := range tests {
		if match := g.Match(v.path); match != v.expected {
			t.Errorf("Incorrect match for '%s' - expected:%v, got:%v", v.path, v.expected, match)
		}
	}
}

func TestGlobMatchPathSlash(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", false},
		{"/alpha", false},
		{"/alpha/pending", true},
		{"/alpha/photos", true},
		{"/alpha/photos/new", false},
		{"/alpha/photos/new/today", false},
		{"/beta", false},
		{"/beta/pending", false},
		{"/beta/photos", false},
	}

	g := NewGlob("/alpha/")

	for _, v := range tests {
		if match := g.Match(v.path); match != v.expected {
			t.Errorf("Incorrect match for '%s' - expected:%v, got:%v", v.path, v.expected, match)
		}
	}
}

func TestGlobMatchPathSlashStar(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", false},
		{"/alpha", false},
		{"/alpha/pending", true},
		{"/alpha/photos", true},
		{"/alpha/photos/new", false},
		{"/alpha/photos/new/today", false},
		{"/beta", false},
		{"/beta/pending", false},
		{"/beta/photos", false},
	}

	g := NewGlob("/alpha/*")

	for _, v := range tests {
		if match := g.Match(v.path); match != v.expected {
			t.Errorf("Incorrect match for '%s' - expected:%v, got:%v", v.path, v.expected, match)
		}
	}
}

func TestGlobMatchPathSlashStarStar(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", false},
		{"/alpha", false},
		{"/alpha/pending", true},
		{"/alpha/photos", true},
		{"/alpha/photos/new", true},
		{"/alpha/photos/new/today", true},
		{"/beta", false},
		{"/beta/pending", false},
		{"/beta/photos", false},
	}

	g := NewGlob("/alpha/**")

	for _, v := range tests {
		if match := g.Match(v.path); match != v.expected {
			t.Errorf("Incorrect match for '%s' - expected:%v, got:%v", v.path, v.expected, match)
		}
	}
}
