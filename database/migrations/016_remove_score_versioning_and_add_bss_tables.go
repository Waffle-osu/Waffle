package migrations

import (
	"database/sql"
)

type MigrationRemoveScoreVersioningAndBssTables struct{}

func (migration MigrationRemoveScoreVersioningAndBssTables) Apply(db *sql.DB) error {
	sql := `
		ALTER TABLE waffle.scores DROP SYSTEM VERSIONING;
		@@@@
		ALTER TABLE waffle.scores ADD COLUMN client_version BIGINT NOT NULL DEFAULT '0' AFTER score_hash
		@@@@
		ALTER TABLE waffle.scores ADD SYSTEM VERSIONING
		@@@@
		CREATE TABLE osu_beatmap_posts (
			beatmapset_id INT NOT NULL,
			subject TEXT NOT NULL,
			message TEXT NOT NULL,
			notify TINYINT NOT NULL,
			complete TINYINT NOT NULL,
			
			PRIMARY KEY(beatmapset_id)
		)
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationRemoveScoreVersioningAndBssTables) Remove(db *sql.DB) error {
	sql :=
		`
		ALTER TABLE waffle.scores DROP SYSTEM VERSIONING;
		@@@@
		ALTER TABLE waffle.scores DROP COLUMN client_version;
		@@@@
		ALTER TABLE waffle.scores ADD SYSTEM VERSIONING
		@@@@
		DROP TABLE osu_beatmap_posts;
	`
	return MigrationHelperRunSplitSql(sql, db)
}
