package options

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
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

func TestHas(t *testing.T) {
	opts := Options{
		{"h", false, "Show Help"},
	}
	if opts.Has("h") != true {
		t.Fatal(`Has("h") should return true but false`)
	}
	if opts.Has("g") != false {
		t.Fatal(`Has("g") should return true but false`)
	}
}

func TestBool(t *testing.T) {
	opts := Options{
		{"h", false, "Show Help"},
	}
	if opts.Bool("h") != false {
		t.Fatal(`Bool("h") should return false but true`)
	}
	if opts.Bool("g") != false {
		t.Fatal(`Bool("g") should return false but true`)
	}
}

func TestIsBool(t *testing.T) {
	opts := Options{
		{"h", false, "Show Help"},
	}
	if opts.IsBool("h") != true {
		t.Fatal(`IsBool("h") should return false but true`)
	}
	if opts.IsBool("g") != false {
		t.Fatal(`IsBool("g") should return false but true`)
	}
}

func TestString(t *testing.T) {
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	if opts.String("foo") != "bar" {
		t.Fatal(`String("foo") should return bar`)
	}
	if opts.String("boo") != "" {
		t.Fatal(`String("boo") should return ""`)
	}
}

func TestParse(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz", "-bar", "foo"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
		{"bar", "baz", "Specify bar"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	if err := opts.Parse(); err != nil {
		t.Fatal(err)
	}
	if opts.String("foo") != "baz" {
		t.Fatal(`String("foo") should return baz`)
	}
	if opts.String("bar") != "foo" {
		t.Fatal(`String("bar") should return foo`)
	}
	if opts.Bool("h") != true {
		t.Fatal(`Bool("h") should return true but false`)
	}
}

func TestParseFail(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz", "-boo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	if err := opts.Parse(); err == nil {
		t.Fatal("Parse Should be fail")
	}
}

func TestParseUnknown(t *testing.T) {
	args := []string {"gotest", "-h", "-"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	if err := opts.Parse(); err == nil {
		t.Fatal(`"-" should be error`)
	}
}

func TestParseDash(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz", "noo", "--", "-boo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	if err := opts.Parse(); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, arg := range Args {
		if arg == "-boo=baz" {
			found = true
		}
	}
	if !found {
		t.Fatal("Dash should keep -boo=baz as an argument")
	}
	found = false
	for _, arg := range Args {
		if arg == "noo" {
			found = true
		}
	}
	if !found {
		t.Fatal(`Args should contain "noo"`)
	}
}

func TestTypeMismatch(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
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

func TestPrintDefaults(t *testing.T) {
	args := []string {"gotest", "-h", "-foo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	opts.Parse()
	var b bytes.Buffer
	temp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fail()
	}
	oldStdout := os.Stdout
	os.Stdout = temp
	opts.PrintDefaults()
	os.Stdout = oldStdout
	output := string(b.Bytes())
	if strings.Contains(output, "-h(false)") {
		t.Fatal(`PrintDefaults should contains -h(false)`)
	}
	if strings.Contains(output, `-foo("baz")`) {
		t.Fatal(`PrintDefaults should contains -h(false)`)
	}
}

func TestUsage(t *testing.T) {
	called := false
	exit = func(n int) {
		called = true
	}
	args := []string {"gotest", "-h", "-foo=baz"}
	opts := Options{
		{"h", false, "Show Help"},
		{"foo", "bar", "Specify foo"},
	}
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = args
	opts.Parse()
	var b bytes.Buffer
	temp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fail()
	}
	oldStdout := os.Stdout
	os.Stdout = temp
	opts.Usage()
	if !called {
		t.Fatal(`Usage wasn't called`)
	}
	os.Stdout = oldStdout
	output := string(b.Bytes())
	if strings.Contains(output, "-h(false)") {
		t.Fatal(`PrintDefaults should contains -h(false)`)
	}
	if strings.Contains(output, `-foo("baz")`) {
		t.Fatal(`PrintDefaults should contains -h(false)`)
	}

}
