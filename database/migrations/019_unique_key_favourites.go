package migrations

import (
	"database/sql"
)

type MigrationAddUniqueKeyBeatmapFavourites struct{}

func (migration MigrationAddUniqueKeyBeatmapFavourites) Apply(db *sql.DB) error {
	sql := `
		ALTER TABLE waffle.beatmap_favourites ADD UNIQUE unique_index (beatmapset_id, user_id)
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddUniqueKeyBeatmapFavourites) Remove(db *sql.DB) error {
	sql :=
		`
	`
	return MigrationHelperRunSplitSql(sql, db)
}
