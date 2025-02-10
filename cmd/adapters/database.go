package adapters

import (
	cockroachdb "proposal-template/pkg/database/cockroachDB"
	"proposal-template/pkg/logger"

	"github.com/golobby/container/v3"
	"gorm.io/gorm"
)
func IoCDatabase() {
	container.Singleton(func() *gorm.DB {
		var (
			logger  logger.ILogger
		)

		err := container.Resolve(&logger)
		if err != nil {
			panic(err)
		}
	
		db, err := cockroachdb.NewCockroachDB(
			cockroachdb.WithLogger(logger),
		)
		if err != nil {
			panic(err)
		}
		
		cockroachdb.CockroachDBGooseMigrate(db, cockroachdb.CockroachDBMigrateFS, "migrations")
		return db
	})
}