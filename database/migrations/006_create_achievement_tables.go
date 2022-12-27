package migrations

import (
	"database/sql"
)

type AchievementTablesStruct struct{}

func (migration AchievementTablesStruct) Apply(db *sql.DB) error {
	creationSql :=
		`
		CREATE TABLE waffle.osu_achievements (
			achievement_id int         unsigned NOT NULL AUTO_INCREMENT,
			name           varchar(64)          NOT NULL,
			image          varchar(64)          NOT NULL,
			
			PRIMARY KEY (achievement_id, name, image)
		) DEFAULT CHARSET=utf8mb4;
@@@@
		CREATE TABLE waffle.osu_achieved_achievements (
			user_achievement_id bigint              NOT NULL AUTO_INCREMENT,
			achievement_id      int             DEFAULT NULL,
			user_id             bigint unsigned DEFAULT NULL,
			
			PRIMARY KEY (user_achievement_id)
		) DEFAULT CHARSET=utf8mb4;
	`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration AchievementTablesStruct) Remove(db *sql.DB) error {
	deletionSql :=
		`
		DROP TABLE waffle.osu_achievements;
		@@@@
		DROP TABLE waffle.osu_achieved_achievements;
	`

	return MigrationHelperRunSplitSql(deletionSql, db)
}
