package utl

import (
	//	"fmt"
	"testing"
)

func TestGetHostname(t *testing.T) {
	_, err := GetHostname()
	if err != nil {
		t.Error(err)
	}
}

func TestGetIPAddrs(t *testing.T) {
	_, err := GetIPAddrs()
	if err != nil {
		t.Error(err)
	}
}
