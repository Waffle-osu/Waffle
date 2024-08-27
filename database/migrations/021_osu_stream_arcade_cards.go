package migrations

import (
	"database/sql"
)

type MigrationAddArcadeCardStructure struct{}

func (migration MigrationAddArcadeCardStructure) Apply(db *sql.DB) error {
	sql := `
		CREATE TABLE waffle.arcade_cards (
			card_id VARCHAR(128) NOT NULL,
			card_pin VARCHAR(128),
			user_id BIGINT
		)
		@@@@
		CREATE TABLE waffle.arcade_link_codes (
			card_id VARCHAR(128) NOT NULL,
			user_id BIGINT NOT NULL,
			link_code VARCHAR(4) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
		@@@@
		CREATE TABLE waffle.arcade_active_tokens (
			card_id VARCHAR(128) NOT NULL,
			token VARCHAR(128) NOT NULL
		)
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddArcadeCardStructure) Remove(db *sql.DB) error {
	sql :=
		`
		DROP TABLE waffle.arcade_cards
		@@@@
		DROP TABLE waffle.arcade_link_codes

	`
	return MigrationHelperRunSplitSql(sql, db)
}
