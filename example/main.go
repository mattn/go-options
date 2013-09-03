package main

import (
	"fmt"
	"github.com/mattn/go-options"
	"os"
)

var opts = options.Options{
	{"h", false, "Show Help"},
	{"verbose", false, "Verbose output"},
	{"prefix", " ", "Prefix of output"},
}

func main() {
	if err := opts.Parse(); err != nil || opts.Bool("h") {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		opts.Usage()
	}

	for i, arg := range options.Args {
		fmt.Print(opts.String("prefix"))
		if opts.Bool("verbose") {
			fmt.Printf("argument %d is %s\n", i+1, arg)
		} else {
			fmt.Println(arg)
		}
	}
}
