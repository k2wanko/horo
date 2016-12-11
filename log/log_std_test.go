//+build !appengine

package log

import (
	"bytes"
	"testing"

	"golang.org/x/net/context"
)

func TestStdLogWrite(t *testing.T) {
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	l := &logger{out: out, err: errOut}
	c := context.Background()

	l.Debugf(c, "Test")
	l.Infof(c, "Test")
	l.Warnf(c, "Test")
	l.Errorf(c, "Test")
	l.Fatalf(c, "Test")

	t.Logf("Stdout:\n%s", out)
	if out, want := out.String(), "[DEBUG] Test\n[INFO] Test\n"; out != want {
		t.Errorf("out = %s; want = %s", out, want)
	}

	t.Logf("Stderr:\n%s", errOut)
	if out, want := errOut.String(), "[WARN] Test\n[ERROR] Test\n[FATAL] Test\n"; out != want {
		t.Errorf("errOut = %s; want = %s", out, want)
	}
}
