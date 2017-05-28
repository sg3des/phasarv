package scenes

import (
	"os"
	"testing"
)

func init() {
	os.Chdir("../../")
}

func TestLoad(t *testing.T) {
	s, err := read("scene00")
	if err != nil {
		t.Fatal(err)
	}
	if s.Name == "" {
		t.Fatal("name is empty")
	}
}
