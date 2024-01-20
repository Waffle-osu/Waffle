package migrations

import (
	"database/sql"
)

type MigrationAddAdditionalBeatmapData struct{}

func (migration MigrationAddAdditionalBeatmapData) Apply(db *sql.DB) error {
	sql := `
		CREATE TABLE osu_bancho_beatmap_playcounts (
			beatmap_id INT NOT NULL,
			beatmapset_id INT NOT NULL,
			passcount INT NOT NULL,
			playcount INT NOT NULL,
			mode TINYINT NOT NULL,

			PRIMARY KEY(beatmap_id, beatmapset_id)
		);
		@@@@
		ALTER TABLE waffle.beatmapsets ADD COLUMN genre_id TINYINT NOT NULL DEFAULT '0' AFTER bpm
		@@@@
		ALTER TABLE waffle.beatmapsets ADD COLUMN language_id TINYINT NOT NULL DEFAULT '0' AFTER genre_id
		@@@@
		ALTER TABLE waffle.beatmapsets ADD COLUMN beatmap_pack varchar(64) NOT NULL DEFAULT '0' AFTER language_id
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddAdditionalBeatmapData) Remove(db *sql.DB) error {
	sql :=
		`
		DROP TABLE waffle.osu_bancho_beatmap_playcounts
		@@@@
		ALTER TABLE waffle.beatmapsets DROP COLUMN genre_id
		@@@@
		ALTER TABLE waffle.beatmapsets DROP COLUMN language_id
		@@@@
		ALTER TABLE waffle.beatmapsets DROP COLUMN beatmap_pack
	`
	return MigrationHelperRunSplitSql(sql, db)
}
