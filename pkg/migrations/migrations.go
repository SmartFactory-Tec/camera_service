package migrations

import "embed"

// go embed does not support embedding files from parent directories
// in lieu of that, this package contains the embedded migrations

//go:embed *.sql
var Migrations embed.FS
