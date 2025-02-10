package cockroachdb


import (
	"embed"
)

//go:embed migrations/*.sql
var CockroachDBMigrateFS embed.FS