package options

import (
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
