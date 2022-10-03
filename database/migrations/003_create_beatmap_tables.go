package migrations

import (
	"database/sql"
)

type CreateBeatmapTablesStruct struct{}

func (migration CreateBeatmapTablesStruct) Apply(db *sql.DB) error {
	creationSql :=
		`
	CREATE TABLE waffle.beatmaps (
		beatmap_id     int          NOT NULL AUTO_INCREMENT,
		beatmapset_id  int          NOT NULL,
		creator_id     bigint       NOT NULL,
		filename       varchar(256) NOT NULL,
		beatmap_md5    varchar(32)  NOT NULL,
		version        varchar(128) NOT NULL,
		total_length   int          NOT NULL,
		drain_time     int          NOT NULL,
		count_objects  int          NOT NULL,
		count_normal   int          NOT NULL,
		count_slider   int          NOT NULL,
		count_spinner  int          NOT NULL,
		diff_hp        tinyint      NOT NULL,
		diff_cs        tinyint      NOT NULL,
		diff_od        tinyint      NOT NULL,
		diff_stars     float        NOT NULL,
		playmode       tinyint      NOT NULL,
		ranking_status tinyint      NOT NULL,
		last_update    datetime     NOT NULL,
		submit_date    datetime     NOT NULL,
		approve_date   datetime     NOT NULL,
		beatmap_source tinyint      NOT NULL COMMENT 'Where the beatmap is from, used for when i may or may not add BSS and might be used for oldsu maps',
		
		PRIMARY KEY (beatmap_id),
		
		KEY checksum_index (beatmap_md5),
		KEY filename_index (filename),
		KEY setid_index    (beatmap_id) /*!80000 INVISIBLE */,
		KEY source_index   (beatmap_source)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.beatmapsets (
		beatmapset_id  int           NOT NULL AUTO_INCREMENT,
		creator_id     bigint        NOT NULL,
		artist         varchar(256)  NOT NULL,
		title          varchar(256)  NOT NULL,
		creator        varchar(256)  NOT NULL,
		source         varchar(256)  NOT NULL,
		tags           varchar(1024) NOT NULL,
		has_video      tinyint       NOT NULL,
		has_storyboard tinyint       NOT NULL,
		bpm            float         NOT NULL,
		
		PRIMARY KEY (beatmapset_id),
		
		KEY title_index   (title) /*!80000 INVISIBLE */,
		KEY artist_index  (title) /*!80000 INVISIBLE */,
		KEY creator_index (creator),
		
		FULLTEXT KEY fulltext_search (artist,title,creator,source,tags)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.beatmap_ratings (
		beatmapset_id int    NOT NULL,
		rating_sum    bigint NOT NULL DEFAULT '0',
		votes         bigint NOT NULL DEFAULT '0',
		
		PRIMARY KEY (beatmapset_id)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.beatmap_ratings_submissions (
		submission_id bigint          NOT NULL AUTO_INCREMENT,
		user_id       bigint unsigned NOT NULL,
		beatmapset_id int             NOT NULL,
		
		PRIMARY KEY (submission_id, user_id, beatmapset_id)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.beatmap_offsets (
		offset_id  bigint NOT NULL AUTO_INCREMENT,
		beatmap_id int    NOT NULL,
		offset     int    NOT NULL DEFAULT '0',
		
		PRIMARY KEY (offset_id, beatmap_id)
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.beatmap_favourites (
		favourite_id  bigint          NOT NULL AUTO_INCREMENT,
		beatmapset_id int             NOT NULL,
		user_id       bigint unsigned NOT NULL,
		
		PRIMARY KEY (favourite_id, user_id, beatmapset_id),
		
		KEY userid_index (user_id),
		
		CONSTRAINT userid_pk FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
@@@@
	CREATE TABLE waffle.beatmap_comments (
		comment_id    bigint          NOT NULL AUTO_INCREMENT,
		user_id       bigint unsigned NOT NULL,
		beatmap_id    int             NOT NULL,
		beatmapset_id int             NOT NULL,
		score_id      bigint unsigned NOT NULL,
		time          bigint          NOT NULL,
		target        tinyint         NOT NULL,
		comment       text            NOT NULL,
		format_string varchar(16)     NOT NULL,
		
		PRIMARY KEY (comment_id),
		
		KEY beatmapid_index                (beatmap_id, target),
		KEY beatmapset_id_index            (beatmapset_id, target) /*!80000 INVISIBLE */,
		KEY score_id_index                 (score_id, target)      /*!80000 INVISIBLE */,
		KEY userid_beatmap_comments_pk_idx (user_id)               /*!80000 INVISIBLE */,
		
		CONSTRAINT userid_beatmapcomments_pk FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
	) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration CreateBeatmapTablesStruct) Remove(db *sql.DB) error {
	deletionSql :=
		`
		DROP TABLE waffle.beatmaps;
		@@@@
		DROP TABLE waffle.beatmapsets;
		@@@@
		DROP TABLE waffle.beatmap_ratings;
		@@@@
		DROP TABLE waffle.beatmap_ratings_submissions;
		@@@@
		DROP TABLE waffle.beatmap_offsets;
		@@@@
		DROP TABLE waffle.beatmap_favourites;
		@@@@
		DROP TABLE waffle.beatmap_comments;
	`

	return MigrationHelperRunSplitSql(deletionSql, db)
}
