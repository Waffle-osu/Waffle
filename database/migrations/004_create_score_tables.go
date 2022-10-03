package migrations

import "database/sql"

type CreateScoreTablesStruct struct{}

func (migration CreateScoreTablesStruct) Apply(db *sql.DB) {
	creationSql :=
		`
		CREATE TABLE waffle.scores (
			score_id         bigint     unsigned NOT NULL AUTO_INCREMENT,
			beatmap_id       int                 NOT NULL,
			beatmapset_id    int                 NOT NULL,
			user_id          bigint     unsigned NOT NULL,
			playmode         tinyint             NOT NULL,
			score            int                 NOT NULL DEFAULT '0',
			max_combo        int                 NOT NULL DEFAULT '0',
			ranking          varchar(2)          NOT NULL,
			hit300           int                 NOT NULL DEFAULT '0',
			hit100           int                 NOT NULL DEFAULT '0',
			hit50            int                 NOT NULL DEFAULT '0',
			hitMiss          int                 NOT NULL DEFAULT '0',
			hitGeki          int                 NOT NULL DEFAULT '0',
			hitKatu          int                 NOT NULL DEFAULT '0',
			enabled_mods     int                 NOT NULL DEFAULT '0',
			perfect          tinyint             NOT NULL DEFAULT '0',
			passed           tinyint             NOT NULL DEFAULT '0',
			date             datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP,
			leaderboard_best tinyint             NOT NULL DEFAULT '0',
			mapset_best      tinyint             NOT NULL DEFAULT '0',
			score_hash       varchar(64)         NOT NULL,
			
			PRIMARY KEY (score_id),
			
			KEY userid_index          (score_id)         /*!80000 INVISIBLE */,
			KEY leaderboardbest_index (leaderboard_best) /*!80000 INVISIBLE */,
			KEY mapsetbest_index      (mapset_best)      /*!80000 INVISIBLE */,
			KEY scorehash_index       (score_hash)       /*!80000 INVISIBLE */,
			KEY beatmapid_index       (beatmap_id),
			KEY userid_fk_idx         (user_id),
			
			CONSTRAINT userid_fk FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
		) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

		CREATE TABLE waffle.failtimes (
			failtime_id bigint unsigned NOT NULL AUTO_INCREMENT,
			failtime    int             NOT NULL,
			beatmap_id  int             NOT NULL,
			score_id    bigint unsigned NOT NULL,
			was_exit    tinyint         NOT NULL,
			
			PRIMARY KEY (failtime_id),
			
			KEY beatmapid_index (beatmap_id)
		) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
	`

	db.Query(creationSql)
}

func (migration CreateScoreTablesStruct) Remove(db *sql.DB) {
	db.Query("DROP TABLE waffle.scores")
	db.Query("DROP TABLE waffle.failtimes")
}
