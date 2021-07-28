package cmd

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	cmd := NewVersionCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(out), "v0.10.1") {
		t.Fatalf("expected \"%s\" got \"%s\"", "v0.10.1", string(out))
	}
}
