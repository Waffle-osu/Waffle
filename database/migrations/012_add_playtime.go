package migrations

import (
	"database/sql"
)

type MigrationAddPlaytime struct{}

func (migration MigrationAddPlaytime) Apply(db *sql.DB) error {
	sql := `
		ALTER TABLE waffle.stats ADD COLUMN playtime BIGINT NOT NULL DEFAULT 0 AFTER replays_watched;
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddPlaytime) Remove(db *sql.DB) error {
	sql :=
		`
	ALTER TABLE waffle.users DROP COLUMN playtime;
	`
	return MigrationHelperRunSplitSql(sql, db)
}
