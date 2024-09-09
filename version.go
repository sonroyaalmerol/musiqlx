//nolint:gochecknoglobals,golint,stylecheck
package musiqlx

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var version string
var Version = strings.TrimSpace(version)

const (
	Name      = "musiqlx"
	NameUpper = "MUSIQLX"
)
