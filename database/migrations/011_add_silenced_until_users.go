package migrations

import "database/sql"

type MigrationAddSilences struct{}

func (migration MigrationAddSilences) Apply(db *sql.DB) error {
	sql :=
		`
		ALTER TABLE waffle.users ADD COLUMN silenced_until BIGINT NOT NULL DEFAULT 0 AFTER joined_at;
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddSilences) Remove(db *sql.DB) error {
	sql :=
		`
		ALTER TABLE waffle.users DROP COLUMN silenced_until;
	`
	return MigrationHelperRunSplitSql(sql, db)
}
