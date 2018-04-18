package help

import "fmt"

func Root(version string) string {
	return fmt.Sprintf(`Token Command Line Interface, version %v`, version)
}
