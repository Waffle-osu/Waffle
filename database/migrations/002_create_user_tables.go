package migrations

import (
	"database/sql"
	"errors"
)

type CreateUserTablesStruct struct{}

func (migration CreateUserTablesStruct) Apply(db *sql.DB) error {
	creationSql :=
		`
	CREATE TABLE waffle.users (
		user_id       BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT,
		username      VARCHAR(32)           NOT NULL,
		password      VARCHAR(64)           NOT NULL,
		country       SMALLINT     UNSIGNED NOT NULL DEFAULT '0',
		banned        TINYINT               NOT NULL DEFAULT '0',
		banned_reason VARCHAR(256)          NOT NULL DEFAULT 'no reason',
		privileges    INT                   NOT NULL DEFAULT '1',
		joined_at     DATETIME                       DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (user_id, username),

		UNIQUE KEY id_UNIQUE       (user_id),
		UNIQUE KEY username_UNIQUE (username),

		KEY user_INDEX (username, user_id)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.stats (
		user_id         bigint  unsigned NOT NULL,
		mode            tinyint          NOT NULL,
		ranked_score    bigint  unsigned NOT NULL DEFAULT '0',
		total_score     bigint  unsigned NOT NULL DEFAULT '0',
		user_level      double           NOT NULL DEFAULT '0',
		accuracy        float            NOT NULL DEFAULT '0',
		playcount       bigint  unsigned NOT NULL DEFAULT '0',
		count_ssh       bigint  unsigned NOT NULL DEFAULT '0',
		count_ss        bigint  unsigned NOT NULL DEFAULT '0',
		count_sh        bigint  unsigned NOT NULL DEFAULT '0',
		count_s         bigint  unsigned NOT NULL DEFAULT '0',
		count_a         bigint  unsigned NOT NULL DEFAULT '0',
		count_b         bigint  unsigned NOT NULL DEFAULT '0',
		count_c         bigint  unsigned NOT NULL DEFAULT '0',
		count_d         bigint  unsigned NOT NULL DEFAULT '0',
		hit300          bigint  unsigned NOT NULL DEFAULT '0',
		hit100          bigint  unsigned NOT NULL DEFAULT '0',
		hit50           bigint  unsigned NOT NULL DEFAULT '0',
		hitMiss         bigint  unsigned NOT NULL DEFAULT '0',
		hitGeki         bigint  unsigned NOT NULL DEFAULT '0',
		hitKatu         bigint  unsigned NOT NULL DEFAULT '0',
		replays_watched bigint  unsigned NOT NULL DEFAULT '0',
	
		PRIMARY KEY (user_id, mode),
		CONSTRAINT userid FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE		
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.friends (
		user_1 bigint unsigned NOT NULL,
		user_2 bigint unsigned NOT NULL,
		
		PRIMARY KEY (user_1, user_2),
		
		KEY index_user1 (user_1) /*!80000 INVISIBLE */,
		KEY index_user2 (user_2),
		
		CONSTRAINT user_id2_FK FOREIGN KEY (user_2) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT user_id_FK FOREIGN KEY (user_1) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.screenshots (
		id            bigint       NOT NULL AUTO_INCREMENT,
		filename      varchar(256) NOT NULL,
		creation_date datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_id       bigint       NOT NULL,
		
		PRIMARY KEY (id)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration CreateUserTablesStruct) Remove(db *sql.DB) error {
	_, err1 := db.Query("DROP TABLE waffle.users")
	_, err2 := db.Query("DROP TABLE waffle.stats")
	_, err3 := db.Query("DROP TABLE waffle.friends")
	_, err4 := db.Query("DROP TABLE waffle.screenshots")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("Dropping failed!")
	}

	return nil
}
