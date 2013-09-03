package options

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	opt := Option{"h", false, "Show Help"}
	opts := Options{&opt}
	if opts.Get("h") != &opt {
		t.Fatal("Get failed")
	}
	if opts.Get("unknown") != nil {
		t.Fatal("Get should return nil for unknown option")
	}
}

func TestBool(t *testing.T) {
	opts := Options{
		{"h", false, "Show Help"},
	}
	if opts.Bool("h") != false {
		t.Fatal("Bool should return false but true")
	}
}

func TestString(t *testing.T) {
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	if opts.String("foo") != "bar" {
		t.Fatal("String should return bar")
	}
}

func TestParse(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "baz", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	opts.Parse()
	if opts.String("foo") != "baz" {
		t.Fatal("String should return baz")
	}
	if opts.Bool("h") != true {
		t.Fatal("Bool should return true but false")
	}
}

func TestTypeMismatch(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "baz", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	opts.Parse()
	if opts.Bool("foo") != false {
		t.Fatal("String for foo should return false")
	}
	if opts.String("h") != "" {
		t.Fatal(`Bool for h should ""`)
	}
}
