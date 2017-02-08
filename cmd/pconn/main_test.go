package main

import (
	"testing"

	"strings"

	"github.com/msbu-tech/go-pconn/cmd/version"
)

func TestVersion(t *testing.T) {
	if strings.Count(version.Version, "") < 5 {
		t.Error("Version is not correct")
	}
}
