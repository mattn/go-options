package options_test

import (
	"fmt"
	"github.com/mattn/go-options"
	"os"
)

// Example 1: options to define
var opts = options.Options{
	{"h", false, "Show Help"},
	{"verbose", false, "Verbose output"},
	{"prefix", " ", "Prefix of output"},
}

func Example() {
	// Example 2: Parse command line arguments. Check -h.
	if err := opts.Parse(); err != nil || opts.Bool("h") {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		opts.Usage()
	}

	for i, arg := range options.Args {
		// Example 3: Use String/Bool to get the values.
		fmt.Print(opts.String("prefix"))
		if opts.Bool("verbose") {
			fmt.Printf("argument %d is %s\n", i+1, arg)
		} else {
			fmt.Println(arg)
		}
	}
}
