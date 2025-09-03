package server

import (
	"strings"
	"testing"
)

func TestMediaDescription(t *testing.T) {
	m := MediaDB{
		Name: "test",
		Type: "ANIME",
	}

	if !strings.Contains(m.Description(), m.Name) {
		t.Errorf("Description does not contain Name")
	}
	if !strings.Contains(m.Description(), m.Type) {
		t.Errorf("Description does not contain Type")
	}
}
