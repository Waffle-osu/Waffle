package migrations

import (
	"database/sql"
)

type MigrationAddDiffCalcTables struct{}

func (migration MigrationAddDiffCalcTables) Apply(db *sql.DB) error {
	sql := `
		CREATE TABLE osu_beatmap_difficulty (
			beatmap_id INT NOT NULL,
			beatmapset_id INT NOT NULL,

			mode TINYINT NOT NULL,

			eyup_stars DOUBLE,

			category_aim_2014    DOUBLE COMMENT 'Catch and Standard',
			category_speed_2014  DOUBLE COMMENT 'Only Standard',
			category_acc_2014    DOUBLE COMMENT 'Everything but catch',
			category_ar_2014     DOUBLE COMMENT 'Catch and Standard',
			category_strain_2014 DOUBLE COMMENT 'Mania & Taiko',
			category_combo_2014  DOUBLE COMMENT 'Catch only',
			category_score_2014  DOUBLE COMMENT 'Mania only',
			category_total_2014  DOUBLE COMMENT 'Final SR',

			category_aim_2016    DOUBLE COMMENT 'Catch and Standard',
			category_speed_2016  DOUBLE COMMENT 'Only Standard',
			category_acc_2016    DOUBLE COMMENT 'Everything but catch',
			category_ar_2016     DOUBLE COMMENT 'Catch and Standard',
			category_strain_2016 DOUBLE COMMENT 'Mania & Taiko',
			category_combo_2016  DOUBLE COMMENT 'Catch only',
			category_score_2016  DOUBLE COMMENT 'Mania only',
			category_total_2016  DOUBLE COMMENT 'Final SR',

			PRIMARY KEY(beatmap_id, mode)
		)
	`
	return MigrationHelperRunSplitSql(sql, db)
}

func (migration MigrationAddDiffCalcTables) Remove(db *sql.DB) error {
	sql :=
		`
		DROP TABLE osu_beatmap_difficulty
	`
	return MigrationHelperRunSplitSql(sql, db)
}
