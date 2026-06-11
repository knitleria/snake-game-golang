// Package version holds the client build version and helpers to compare
// dotted numeric versions (e.g. "v1.2.3"). The Version variable is injected at
// build time via -ldflags -X and defaults to "dev" for local builds.
package version

import (
	"strconv"
	"strings"
)

// Version is the client build version, injected at build time via:
//
//	-ldflags "-X 'snake_golang/internal/version.Version=v1.2.3'"
//
// Local/dev builds keep the default "dev", which IsNumeric reports as false.
var Version = "dev"

// IsNumeric reports whether v looks like a dotted numeric version (e.g.
// "1.2.3" or "v1.2.3"). Values like "dev" or "" are not numeric.
func IsNumeric(v string) bool {
	parts := parse(v)
	return parts != nil
}

// Compare compares two dotted numeric versions and returns:
//
//	-1 if a < b, 0 if a == b, 1 if a > b
//
// A leading "v"/"V" is ignored. Non-numeric values (e.g. "dev", "") are
// treated as the lowest possible version (0.0.0...).
func Compare(a, b string) int {
	pa := parse(a)
	pb := parse(b)
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		var x, y int
		if i < len(pa) {
			x = pa[i]
		}
		if i < len(pb) {
			y = pb[i]
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return 0
}

// AtLeast reports whether version v is greater than or equal to min.
// An empty min disables the check and always returns true.
func AtLeast(v, min string) bool {
	if strings.TrimSpace(min) == "" {
		return true
	}
	return Compare(v, min) >= 0
}

// parse splits a dotted numeric version into its integer components. It returns
// nil when the value is not a valid dotted numeric version.
func parse(v string) []int {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	if v == "" {
		return nil
	}
	fields := strings.Split(v, ".")
	out := make([]int, 0, len(fields))
	for _, f := range fields {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil || n < 0 {
			return nil
		}
		out = append(out, n)
	}
	return out
}
