package migrations

import (
	"database/sql"
)

type IrcAndUpdaterTablesStruct struct{}

func (migration IrcAndUpdaterTablesStruct) Apply(db *sql.DB) error {
	creationSql :=
		`
		CREATE TABLE waffle.irc_log (
			message_id bigint      unsigned NOT NULL AUTO_INCREMENT,
			sender     bigint      unsigned NOT NULL,
			target     varchar(64)          NOT NULL,
			message    text                 NOT NULL,
			date       datetime             NOT NULL DEFAULT CURRENT_TIMESTAMP,
			
			PRIMARY KEY (message_id)
		) DEFAULT CHARSET=utf8mb4;
@@@@
		CREATE TABLE waffle.updater_items (
			item_id         bigint       NOT NULL AUTO_INCREMENT,
			server_filename varchar(256) NOT NULL,
			client_filename varchar(256) NOT NULL,
			file_hash       varchar(64)  NOT NULL,
			item_name       varchar(256) NOT NULL DEFAULT '',
			item_action     varchar(8)            DEFAULT 'none',
			
			PRIMARY KEY (item_id, server_filename, client_filename)
		) DEFAULT CHARSET=utf8mb4;
@@@@
		INSERT INTO waffle.updater_items (server_filename, client_filename, file_hash, item_name, item_action) VALUES
			('osu!.exe'                   , 'osu!.exe'                   , '3623d9f7c693b786564e2d61b1c43af9', 'client_debug', 'none'),
			('avcodec-51.dll'             , 'avcodec-51.dll'             , 'b22bf1e4ecd4be3d909dc68ccab74eec', 'client_ddls',  'none'),
			('avformat-52.dll'            , 'avformat-52.dll'            , '2e7a800133625f827cf46aa0bb1af800', 'client_ddls',  'none'),
			('avutil-49.dll'              , 'avutil-49.dll'              , 'c870147dff89c95c81f8fbdfbc6344ac', 'client_ddls',  'none'),
			('bass.dll'                   , 'bass.dll'                   , 'bbfc7d855252b0211875769bbf667bcd', 'client_ddls',  'none'),
			('bass_fx.dll'                , 'bass_fx.dll'                , 'f9ffe0a23a32b79653e31330764c4231', 'client_ddls',  'none'),
			('d3dx9_31.dll'               , 'd3dx9_31.dll'               , '797e24743937d67d69f28f2cf5052ee8', 'client_ddls',  'none'),
			('Microsoft.Ink.dll'          , 'Microsoft.Ink.dll'          , 'a02ee61542caae25f8a44c9428d30247', 'client_ddls',  'none'),
			('Microsoft.Xna.Framework.dll', 'Microsoft.Xna.Framework.dll', '45a786658d3f69717652fed471d03ee0', 'client_ddls',  'none'),
			('osu!common.dll'             , 'osu!common.dll'             , '820817b776374a0adcbb7fa2a7ca74f2', 'client_ddls',  'none'),
			('osu.dll.zip'                , 'osu.dll'                    , '599c14bcfc9c43b88d70d1a9388b33b7', 'client_ddls',  'zip' ),
			('OsuP2P.dll'                 , 'OsuP2P.dll'                 , '2342bfd835e2e487d040d8ba62eb1a72', 'client_ddls',  'none'),
			('pthreadGC2.dll'             , 'pthreadGC2.dll'             , 'ce931021e18f385f519e945a8a10548e', 'client_ddls',  'none'),
			('x3daudio1_1.dll'            , 'x3daudio1_1.dll'            , '121b131eaa369d8f58dacc5c39a77d80', 'client_ddls',  'none');
	`

	return MigrationHelperRunSplitSql(creationSql, db)
}

func (migration IrcAndUpdaterTablesStruct) Remove(db *sql.DB) error {
	deletionSql :=
		`
		DROP TABLE waffle.irc_log;
		@@@@
		DROP TABLE waffle.updater_items;	
	`

	return MigrationHelperRunSplitSql(deletionSql, db)
}
