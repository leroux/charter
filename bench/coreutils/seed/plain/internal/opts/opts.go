// Package opts splits a tool's arguments into options and operands, the
// way GNU's getopt does for the options the tools know.
package opts

import (
	"fmt"
	"strings"
)

// Opt is one option as given: its short name and, for an option that
// takes one, its value.
type Opt struct {
	Name  byte
	Value string
}

// Spec is what a tool knows: its short options without a value, those
// with one, and its long options by the short name each stands for.
type Spec struct {
	Flags  string
	Valued string
	Long   map[string]byte
}

// Parse splits args into options, in the order given, and operands.
// Options may come anywhere before "--"; short ones may be clustered, and
// a valued one takes the rest of its cluster or the next argument. "-"
// alone is an operand. An option the tool does not know, or a value
// missing, is an error in GNU's words.
func Parse(args []string, spec Spec) ([]Opt, []string, error) {
	var opts []Opt
	var operands []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--":
			return opts, append(operands, args[i+1:]...), nil
		case strings.HasPrefix(a, "--"):
			name, value, hasValue := strings.Cut(a[2:], "=")
			short, ok := spec.Long[name]
			if !ok {
				return nil, nil, fmt.Errorf("unrecognized option '%s'", a)
			}
			if strings.IndexByte(spec.Valued, short) >= 0 {
				if !hasValue {
					if i+1 >= len(args) {
						return nil, nil, fmt.Errorf("option '--%s' requires an argument", name)
					}
					i++
					value = args[i]
				}
			} else if hasValue {
				return nil, nil, fmt.Errorf("option '--%s' doesn't allow an argument", name)
			}
			opts = append(opts, Opt{Name: short, Value: value})
		case len(a) > 1 && a[0] == '-':
			for j := 1; j < len(a); j++ {
				c := a[j]
				switch {
				case strings.IndexByte(spec.Valued, c) >= 0:
					value := a[j+1:]
					if value == "" {
						if i+1 >= len(args) {
							return nil, nil, fmt.Errorf("option requires an argument -- '%c'", c)
						}
						i++
						value = args[i]
					}
					opts = append(opts, Opt{Name: c, Value: value})
					j = len(a)
				case strings.IndexByte(spec.Flags, c) >= 0:
					opts = append(opts, Opt{Name: c})
				default:
					return nil, nil, fmt.Errorf("invalid option -- '%c'", c)
				}
			}
		default:
			operands = append(operands, a)
		}
	}
	return opts, operands, nil
}
