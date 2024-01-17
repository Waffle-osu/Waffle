package migrations

import (
	"database/sql"
)

type MigrationDropSiteTokens struct{}

func (migration MigrationDropSiteTokens) Apply(db *sql.DB) error {
	sql := `
		DROP TABLE site_tokens;
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationDropSiteTokens) Remove(db *sql.DB) error {
	//never used, fuck it
	return nil
}
