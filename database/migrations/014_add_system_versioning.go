package migrations

import (
	"database/sql"
)

type MigrationAddSystemVersioning struct{}

func (migration MigrationAddSystemVersioning) Apply(db *sql.DB) error {
	sql := `
		ALTER TABLE waffle.scores ADD SYSTEM VERSIONING
		@@@@
		ALTER TABLE waffle.stats ADD SYSTEM VERSIONING
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddSystemVersioning) Remove(db *sql.DB) error {
	sql :=
		`
		ALTER TABLE waffle.scores DROP SYSTEM VERSIONING
		@@@@
		ALTER TABLE waffle.stats DROp SYSTEM VERSIONING
	`
	return MigrationHelperRunSplitSql(sql, db)
}
