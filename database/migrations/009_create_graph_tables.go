package migrations

import (
	"database/sql"
)

type CreateHistoricalTablesStruct struct{}

func (migration CreateHistoricalTablesStruct) Apply(db *sql.DB) error {
	creationSql :=
		`
		CREATE TABLE waffle.osu_historical_stats (
			id              BIGINT            NOT NULL AUTO_INCREMENT,
			mode            tinyint           NOT NULL,
			ranked_score    bigint   unsigned NOT NULL DEFAULT '0',
			total_score     bigint   unsigned NOT NULL DEFAULT '0',
			user_level      double            NOT NULL DEFAULT '0',
			accuracy        float             NOT NULL DEFAULT '0',
			playcount       bigint   unsigned NOT NULL DEFAULT '0',
			count_ssh       bigint   unsigned NOT NULL DEFAULT '0',
			count_ss        bigint   unsigned NOT NULL DEFAULT '0',
			count_sh        bigint   unsigned NOT NULL DEFAULT '0',
			count_s         bigint   unsigned NOT NULL DEFAULT '0',
			count_a         bigint   unsigned NOT NULL DEFAULT '0',
			count_b         bigint   unsigned NOT NULL DEFAULT '0',
			count_c         bigint   unsigned NOT NULL DEFAULT '0',
			count_d         bigint   unsigned NOT NULL DEFAULT '0',
			hit300          bigint   unsigned NOT NULL DEFAULT '0',
			hit100          bigint   unsigned NOT NULL DEFAULT '0',
			hit50           bigint   unsigned NOT NULL DEFAULT '0',
			hitMiss         bigint   unsigned NOT NULL DEFAULT '0',
			hitGeki         bigint   unsigned NOT NULL DEFAULT '0',
			hitKatu         bigint   unsigned NOT NULL DEFAULT '0',
			replays_watched bigint   unsigned NOT NULL DEFAULT '0',
			current_rank    bigint   unsigned NOT NULL DEFAULT '1',
			date            DATETIME          NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(id)
		);
	`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration CreateHistoricalTablesStruct) Remove(db *sql.DB) error {
	deletionSql :=
		`
		DROP TABLE waffle.osu_historical_stats;
	`
	return MigrationHelperRunSplitSql(deletionSql, db)
}
