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

func (options Options) Parse() error {
	hasDash := false
	nArgs := len(os.Args)

	getDefaults(options)
	for n := 1; n < nArgs; n++ {
		arg := os.Args[n]
		if arg == "--" {
			hasDash = true
			continue
		}

		flag, value := "", ""
		if strings.HasPrefix(arg, "-") {
			tokens := strings.SplitN(arg, "=", 2)
			if len(tokens) ==1 {
				flag = tokens[0]
				if n < nArgs-1 && options.Has(flag) && !options.IsBool(flag) {
					value = os.Args[n+1]
					n++
				}
			} else {
				flag = tokens[0]
				value = tokens[1]
			}
		}

		if !hasDash {
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
		} else {
			Args = append(Args, arg)
		}
	}
	return nil
}

var exit = os.Exit

func (options Options) Usage() {
	fmt.Printf("Usage: %s [options] [--] [args]\n", os.Args[0])
	options.PrintDefaults()
	exit(1)
}

func (options Options) PrintDefaults() {
	getDefaults(options)
	for _, option := range options {
		if _, ok := option.Value.(string); ok {
			fmt.Printf("  -%s=%s: %s\n", option.Flag, defaults[option.Flag], option.Description)
		} else {
			fmt.Printf("  -%s(%v): %s\n", option.Flag, defaults[option.Flag], option.Description)
		}
	}
}

func (options Options) Has(flag string) bool {
	for _, option := range options {
		if option.Flag == flag {
			return true
		}
	}
	return false
}

func (options Options) Get(flag string) *Option {
	for _, option := range options {
		if option.Flag == flag {
			return option
		}
	}
	return nil
}

func (options Options) String(flag string) string {
	for _, option := range options {
		if option.Flag == flag {
			s, _ := option.Value.(string)
			return s
		}
	}
	return ""
}

func (options Options) Bool(flag string) bool {
	for _, option := range options {
		if option.Flag == flag {
			b, _ := option.Value.(bool)
			return b
		}
	}
	return false
}

func (options Options) IsBool(flag string) bool {
	for _, option := range options {
		if option.Flag == flag {
			_, ok := option.Value.(bool)
			return ok
		}
	}
	return false
}
