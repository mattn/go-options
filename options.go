//Easily way to get command line flags
package options

import (
	"fmt"
	"os"
	"strings"
)

type Option struct {
	Flag        string
	Value       interface{}
	Description string
}

type Options []*Option

//Args is parsed arguments except options.
var Args []string
var defaults map[string]interface{}

func getDefaults(options Options) {
	if len(defaults) == 0 {
		defaults = make(map[string]interface{})
		for _, option := range options {
			defaults[option.Flag] = option.Value
		}
	}
}

//Parse command line arguments.
//If arguments contains --, following arguments doesn't treat options.
func Parse(options Options) error {
	nArgs := len(os.Args)

	getDefaults(options)
	for n := 1; n < nArgs; n++ {
		arg := os.Args[n]
		if arg == "--" {
			Args = append(Args, os.Args[n+1:]...)
			break
		}

		flag, value := "", ""
		if strings.HasPrefix(arg, "-") {
			tokens := strings.SplitN(arg, "=", 2)
			if len(tokens) ==1 {
				flag = tokens[0]
				if n < nArgs-1 && options.Has(flag[1:]) && !options.IsBool(flag[1:]) {
					value = os.Args[n+1]
					n++
				}
			} else {
				flag = tokens[0]
				value = tokens[1]
			}
		}

		if len(flag) > 0 && flag[0] == '-' {
			option := options.Get(flag[1:])
			if option == nil {
				return fmt.Errorf("Invalid option: '%v'\n", flag)
			} else {
				if _, ok := option.Value.(string); ok {
					option.Value = value
				} else {
					option.Value = true
				}
			}
		} else {
			Args = append(Args, arg)
		}
	}
	return nil
}

//Syntax sugar of options.Parse().
func (options Options) Parse() error {
	return Parse(options)
}

var exit = os.Exit

//Usage shows command line usage and default values. And exit.
func (options Options) Usage() {
	fmt.Printf("Usage: %s [options] [--] [args]\n", os.Args[0])
	options.PrintDefaults()
	exit(1)
}

//PrintDefaults shows default values.
func (options Options) PrintDefaults() {
	getDefaults(options)
	for _, option := range options {
		if _, ok := option.Value.(string); ok {
			fmt.Printf("  -%s=%q: %s\n", option.Flag, defaults[option.Flag], option.Description)
		} else {
			fmt.Printf("  -%s(%v): %s\n", option.Flag, defaults[option.Flag], option.Description)
		}
	}
}

//Has return true if options has the flag.
func (options Options) Has(flag string) bool {
	for _, option := range options {
		if option.Flag == flag {
			return true
		}
	}
	return false
}

//Get return Option struct of specified flag.
func (options Options) Get(flag string) *Option {
	for _, option := range options {
		if option.Flag == flag {
			return option
		}
	}
	return nil
}

//String return the string value of the flag.
func (options Options) String(flag string) string {
	for _, option := range options {
		if option.Flag == flag {
			s, _ := option.Value.(string)
			return s
		}
	}
	return ""
}

//Bool return the boolean value of the flag.
func (options Options) Bool(flag string) bool {
	for _, option := range options {
		if option.Flag == flag {
			b, _ := option.Value.(bool)
			return b
		}
	}
	return false
}

//IsBool return whether the flag is boolean flag or not.
func (options Options) IsBool(flag string) bool {
	for _, option := range options {
		if option.Flag == flag {
			_, ok := option.Value.(bool)
			return ok
		}
	}
	return false
}
