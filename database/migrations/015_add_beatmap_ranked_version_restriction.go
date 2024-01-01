package migrations

import (
	"database/sql"
)

type MigrationAddBeatmapRankedVersionRestriction struct{}

func (migration MigrationAddBeatmapRankedVersionRestriction) Apply(db *sql.DB) error {
	sql := `
		ALTER TABLE waffle.beatmaps ADD COLUMN status_valid_from_version BIGINT NOT NULL DEFAULT '0' AFTER beatmap_source
		@@@@
		ALTER TABLE waffle.beatmaps ADD COLUMN status_valid_to_version BIGINT NOT NULL DEFAULT '0' AFTER status_valid_from_version
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddBeatmapRankedVersionRestriction) Remove(db *sql.DB) error {
	sql :=
		`
		ALTER TABLE waffle.beatmaps DROP COLUMN status_valid_from_version
		@@@@
		ALTER TABLE waffle.beatmaps DROP COLUMN status_valid_to_version
	`
	return MigrationHelperRunSplitSql(sql, db)
}
