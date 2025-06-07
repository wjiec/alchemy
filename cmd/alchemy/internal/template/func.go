package template

import "github.com/spf13/pflag"

// flag creates a template function that retrieves flag values from a FlagSet.
func flag(fs *pflag.FlagSet) func(string) string {
	return func(name string) string {
		var value string
		fs.Visit(func(f *pflag.Flag) {
			if f.Name == name {
				value = f.Value.String()
			}
		})
		return value
	}
}
