package migrations

import (
	"database/sql"
)

type MigrationAddMatchLogs struct{}

func (migration MigrationAddMatchLogs) Apply(db *sql.DB) error {
	sql := `
		CREATE TABLE osu_match_history (
			event_id           BIGINT      UNSIGNED NOT NULL AUTO_INCREMENt
			match_id           VARCHAR(64)          NOT NULL,
			event_type         TINYINT              NOT NULL,
			event_initiator_id BIGINT      UNSIGNED NOT NULL,
			extra_info         TEXT                 NOT NULL,

			PRIMARY KEY (event_id),

			KEY match_id_INDEX (match_id)
		)
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddMatchLogs) Remove(db *sql.DB) error {
	sql :=
		`
		DROP TABLE osu_match_history;
	`
	return MigrationHelperRunSplitSql(sql, db)
}
