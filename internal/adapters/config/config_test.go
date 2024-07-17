package config

import (
	"fmt"
	"testing"
)

// test that the config file is being read correctly
func TestNew(t *testing.T) {
	config, err := Crear()
	if err != nil {
		t.Errorf("error reading config file: %v", err)
	}
	fmt.Println(config)
}
